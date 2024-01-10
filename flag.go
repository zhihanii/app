package app

import (
	"flag"
	"fmt"
	"github.com/spf13/pflag"
	"strings"
)

const versionFlagName = "version"

func WordSepNormalizeFunc(f *pflag.FlagSet, name string) pflag.NormalizedName {
	if strings.Contains(name, "_") {
		return pflag.NormalizedName(strings.ReplaceAll(name, "_", "-"))
	}
	return pflag.NormalizedName(name)
}

func InitFlags(flags *pflag.FlagSet) {
	flags.SetNormalizeFunc(WordSepNormalizeFunc)
	flags.AddGoFlagSet(flag.CommandLine)
}

func PrintFlags(flags *pflag.FlagSet) {
	flags.VisitAll(func(flag *pflag.Flag) {
		//log.Debugf
	})
}

func AddGlobalFlags(fs *pflag.FlagSet, name string) {
	fs.BoolP("help", "h", false, fmt.Sprintf("help for %s", name))
}
