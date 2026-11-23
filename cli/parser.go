package cli

import (
	"strings"
)

type cli interface {
	parse(osArgs []string)
}

type SboAdd struct {
	Arg string
	Dflag bool
	Cflag bool
	Sflag bool
	Clean bool
}

func (c *SboAdd) parse(osArgs []string) {
	for _, arg := range osArgs {
		if !strings.HasPrefix(arg, "-") {
			c.Arg = arg
			continue
		} else {
			for _, char := range arg[1:] {
				switch char {
					case 'd':
						c.Dflag = true
					case 'c':
						c.Cflag = true
					case 'S':
						c.Sflag = true
					case 'C':
						c.Clean = true
				}
			}
		}
	}
}

func Run(c cli, osArgs []string, op func() error) error {
	c.parse(osArgs)
	return op()
}
