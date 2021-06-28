package vals_test

import (
	"flag"
	"strings"
	"testing"

	"github.com/mailund/cli/internal/vals"
)

func simplifyUsage(x string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(x)), " ")
}

func TestValsFlags(t *testing.T) {
	var (
		i vals.IntValue    = 42
		s vals.StringValue = "foo"
		f                  = flag.NewFlagSet("test", flag.ContinueOnError)
	)

	f.Var(&i, "i", "an integer")
	f.Var(&s, "s", "a string")

	builder := new(strings.Builder)
	f.SetOutput(builder)
	f.PrintDefaults()

	usage := simplifyUsage(builder.String())
	expected := strings.Join(
		[]string{
			"-i value an integer (default 42)",
			"-s value a string (default \"foo\")",
		}, " ")

	if usage != expected {
		t.Errorf("Unexpected usage string:\n'%s'\n", usage)
	}

	err := f.Parse([]string{"-i=13", "-s", "bar"})

	if err != nil {
		t.Fatal("Didn't expect parsing to return an error")
	}

	if i != 13 {
		t.Errorf("Expected i to have the value 13, it has %d", i)
	}

	if s != "bar" {
		t.Errorf("Expected s to have the value bar, it has %s", s)
	}
}
