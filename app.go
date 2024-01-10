package app

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zhihanii/app/flag"
)

type App struct {
	name      string
	shortDesc string
	longDesc  string
	options   CmdLineOptions
	runFunc   RunFunc
	silence   bool
	noVersion bool
	noConfig  bool
	commands  []*Command
	args      cobra.PositionalArgs
	cmd       *cobra.Command
}

type RunFunc func(name string) error

type Option func(*App)

func WithOptions(options CmdLineOptions) Option {
	return func(a *App) {
		a.options = options
	}
}

func WithRunFunc(runFunc RunFunc) Option {
	return func(a *App) {
		a.runFunc = runFunc
	}
}

func WithDescription(desc string) Option {
	return func(a *App) {
		a.longDesc = desc
	}
}

func WithSilence() Option {
	return func(a *App) {
		a.silence = true
	}
}

func WithNoVersion() Option {
	return func(a *App) {
		a.noVersion = true
	}
}

func WithNoConfig() Option {
	return func(a *App) {
		a.noConfig = true
	}
}

func WithArgs(args cobra.PositionalArgs) Option {
	return func(a *App) {
		a.args = args
	}
}

func WithDefaultArgs() Option {
	return func(a *App) {
		a.args = func(cmd *cobra.Command, args []string) error {
			for _, arg := range args {
				if len(arg) > 0 {
					return fmt.Errorf("%q does not take any arguments, got %q", cmd.CommandPath(), args)
				}
			}
			return nil
		}
	}
}

func WithCommands(commands []*Command) Option {
	return func(a *App) {
		a.commands = commands
	}
}

func New(name string, shortDesc string, opts ...Option) *App {
	a := &App{
		name:      name,
		shortDesc: shortDesc,
	}

	for _, o := range opts {
		o(a)
	}

	a.buildCommand()

	return a
}

func (a *App) buildCommand() {
	cmd := &cobra.Command{
		Use:           a.name,
		Short:         a.shortDesc,
		Long:          a.longDesc,
		SilenceUsage:  true,
		SilenceErrors: true,
		Args:          a.args,
	}
	cmd.SetOut(os.Stdout)
	cmd.SetErr(os.Stderr)
	cmd.Flags().SortFlags = true
	InitFlags(cmd.Flags())

	if len(a.commands) > 0 {
		for _, command := range a.commands {
			cmd.AddCommand(command.cobraCommand())
		}
		cmd.SetHelpCommand(helpCommand(a.name))
	}
	if a.runFunc != nil {
		cmd.RunE = a.runCommand
	}

	var namedFlagSets flag.NamedFlagSets
	if a.options != nil {
		namedFlagSets = a.options.Flags()
		fs := cmd.Flags()
		for _, f := range namedFlagSets.FlagSets {
			fs.AddFlagSet(f)
		}
	}

	// if !a.noVersion {
	// 	namedFlagSets.FlagSet("global").AddFlag(pflag.Lookup(versionFlagName))
	// }
	if !a.noConfig {
		addConfigFlag(a.name, namedFlagSets.FlagSet("global"))
	}
	AddGlobalFlags(namedFlagSets.FlagSet("global"), cmd.Name())
	cmd.Flags().AddFlagSet(namedFlagSets.FlagSet("global"))

	//addCmdTemplate(cmd, namedFlagSets)
	a.cmd = cmd
}

func (a *App) runCommand(cmd *cobra.Command, args []string) error {
	//PrintFlags(cmd.Flags())
	if !a.noVersion {

	}

	if !a.noConfig {
		//建立Flag.Name => Flag的映射
		if err := viper.BindPFlags(cmd.Flags()); err != nil {
			return err
		}
		//解析配置
		if err := viper.Unmarshal(a.options); err != nil {
			return err
		}
	}

	if !a.silence {
		//log.Infof
		if !a.noVersion {

		}
		if !a.noConfig {
			// log.Printf("config file used: %s\n", viper.ConfigFileUsed())
		}
	}

	if a.options != nil {
		if err := a.applyOptionRules(); err != nil {
			return err
		}
	}

	if a.runFunc != nil {
		return a.runFunc(a.name)
	}

	return nil
}

func (a *App) applyOptionRules() error {
	if completableOptions, ok := a.options.(CompletableOptions); ok {
		if err := completableOptions.Complete(); err != nil {
			return err
		}
	}
	if errs := a.options.Validate(); len(errs) != 0 {

	}
	//if printableOptions, ok := a.options.(PrintableOptions); ok && !a.silence {
	//	//log.Infof
	//}
	return nil
}

func (a *App) Run() {
	if err := a.cmd.Execute(); err != nil {
		fmt.Printf("%v %v\n", color.RedString("Error:"), err)
		os.Exit(1)
	}
}

func addCmdTemplate(cmd *cobra.Command, namedFlagSets flag.NamedFlagSets) {
	usageFmt := "Usage:\n  %s\n"
	//cols, _, _ := term.TerminalSize(cmd.OutOrStdout())
	cmd.SetUsageFunc(func(cmd *cobra.Command) error {
		fmt.Fprintf(cmd.OutOrStderr(), usageFmt, cmd.UseLine())
		//cliflag.PrintSections(cmd.OutOrStderr(), namedFlagSets, cols)

		return nil
	})
	cmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		fmt.Fprintf(cmd.OutOrStdout(), "%s\n\n"+usageFmt, cmd.Long, cmd.UseLine())
		//cliflag.PrintSections(cmd.OutOrStdout(), namedFlagSets, cols)
	})
}
