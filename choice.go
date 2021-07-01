package cli

import (
	"strings"

	"github.com/mailund/cli/interfaces"
)

// Choice is a type for handling arguments that must be from a small set of optoins
type Choice struct {
	Choice  string   // The default/provided choice
	Options []string // The options to choose from
}

// Set implements the FlagValue/PosValue interface
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

// String implements the FlagValue interface
func (c *Choice) String() string {
	return c.Choice
}

// ArgumentDescription implements the ArgumentDescription protocol
func (c *Choice) ArgumentDescription(flag bool, descr string) string {
	if flag {
		// Don't modify the flag description (we handle it on the value)
		return descr
	}

	return descr + " (choose from " + "{" + strings.Join(c.Options, ",") + "})"
}

// FlagValueDescription implements the FlagValueDescription protocol
func (c *Choice) FlagValueDescription() string {
	return "{" + strings.Join(c.Options, ",") + "}"
}
