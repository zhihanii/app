package flag

import "github.com/spf13/pflag"

type NamedFlagSets struct {
	Order    []string
	FlagSets map[string]*pflag.FlagSet
}

func (n *NamedFlagSets) FlagSet(name string) *pflag.FlagSet {
	if n.FlagSets == nil {
		n.FlagSets = make(map[string]*pflag.FlagSet)
	}
	if _, ok := n.FlagSets[name]; !ok {
		n.FlagSets[name] = pflag.NewFlagSet(name, pflag.ExitOnError)
		n.Order = append(n.Order, name)
	}
	return n.FlagSets[name]
}
