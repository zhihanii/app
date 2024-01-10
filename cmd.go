package app

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

type Command struct {
	use   string
	short string
	// options  CmdLineOptions
	commands []*Command
	runFunc  CommandRunFunc
}

type CommandRunFunc func(args []string) error

type CommandOption func(*Command)

func WithCommandRunFunc(runFunc CommandRunFunc) CommandOption {
	return func(c *Command) {
		c.runFunc = runFunc
	}
}

func NewCommand(use string, short string, opts ...CommandOption) *Command {
	c := &Command{
		use:   use,
		short: short,
	}

	for _, o := range opts {
		o(c)
	}

	return c
}

func (c *Command) AddCommand(cmd *Command) {
	c.commands = append(c.commands, cmd)
}

func (c *Command) AddCommands(cmds ...*Command) {
	c.commands = append(c.commands, cmds...)
}

func (c *Command) cobraCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   c.use,
		Short: c.short,
	}
	cmd.SetOut(os.Stdout)
	cmd.SetErr(os.Stderr)
	cmd.Flags().SortFlags = false
	if len(c.commands) > 0 {
		for _, command := range c.commands {
			cmd.AddCommand(command.cobraCommand())
		}
	}
	if c.runFunc != nil {
		cmd.Run = c.run
	}
	// if c.options != nil {
	// 	for _, f := range c.options.Flags().FlagSets {
	// 		cmd.Flags().AddFlagSet(f)
	// 	}
	// }
	addHelpCommandFlag(c.use, cmd.Flags())
	return cmd
}

func (c *Command) run(cmd *cobra.Command, args []string) {
	if c.runFunc != nil {
		if err := c.runFunc(args); err != nil {
			fmt.Printf("%v %v\n", color.RedString("Error:"), err)
			os.Exit(1)
		}
	}
}
