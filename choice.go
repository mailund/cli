package cli

import (
	"strings"

	"github.com/mailund/cli/interfaces"
)

type Choice struct {
	Choice  string
	Options []string
}

func (c *Choice) Set(x string) error {
	for _, v := range c.Options {
		if v == x {
			c.Choice = v
			return nil
		}
	}

	return interfaces.ParseErrorf(
		"%s is not a valid choice, must be in %s", x, "{"+strings.Join(c.Options, ",")+"}")
}

func (c *Choice) String() string {
	return c.Choice
}

func (c *Choice) ArgumentDescription(flag bool, descr string) string {
	if flag {
		// Don't modify the flag description (we handle it on the value)
		return descr
	}

	return descr + " (choose from " + "{" + strings.Join(c.Options, ",") + "})"
}

func (c *Choice) FlagValueDescription() string {
	return "{" + strings.Join(c.Options, ",") + "}"
}
