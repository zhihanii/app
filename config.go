package app

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const configFlagName = "config"

var cfgFile string

func init() {
	pflag.StringVarP(&cfgFile, "config", "c", cfgFile, "Read configuration from specified `FILE`, "+
		"support YAML.")
}

func addConfigFlag(name string, fs *pflag.FlagSet) {
	fs.AddFlag(pflag.Lookup(configFlagName))

	viper.AutomaticEnv()
	viper.SetEnvPrefix(strings.Replace(strings.ToUpper(name), "-", "_", -1))
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	cobra.OnInitialize(func() {
		if cfgFile != "" {
			viper.SetConfigFile(cfgFile)
		} else { //使用默认配置文件
			viper.AddConfigPath(".")

			if names := strings.Split(name, "-"); len(names) > 1 {
				viper.AddConfigPath(filepath.Join(homeDir(), "."+names[0]))
				viper.AddConfigPath(filepath.Join("/etc", names[0]))
			}

			viper.SetConfigName(name)
		}

		if err := viper.ReadInConfig(); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error: failed to read configuration file(%s): %v\n", cfgFile, err)
			os.Exit(1)
		}
	})
}

func homeDir() string {
	return os.Getenv("HOME")
}
