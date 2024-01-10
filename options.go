package app

import (
	"github.com/zhihanii/app/flag"
)

type CmdLineOptions interface {
	Flags() flag.NamedFlagSets
	Validate() []error
}

type ConfigurableOptions interface {
	ApplyFlags() []error
}

type CompletableOptions interface {
	Complete() error
}

type PrintableOptions interface {
	String() string
}
