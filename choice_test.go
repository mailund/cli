package cli_test

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/mailund/cli"
	"github.com/mailund/cli/internal/failure"
)

type Args struct {
	A cli.Choice `flag:"" short:"a" descr:"optional choice"`
	B cli.Choice `pos:"b" descr:"mandatory choice"`
}

var args = &Args{
	A: cli.Choice{"A", []string{"A", "B", "C"}},
	B: cli.Choice{"B", []string{"A", "B", "C"}},
}

func Init() interface{} { return args }

var cmd = cli.NewCommand(
	cli.CommandSpec{
		Name: "choices",
		Long: "Demonstration of the difficult task of making choices.",
		Init: Init,
	})

func TestUsage(t *testing.T) {
	builder := new(strings.Builder)
	cmd.SetOutput(builder)

	cmd.Usage()

	expected := `Usage: choices [flags] b

	Demonstration of the difficult task of making choices.

	Flags:
	-h,--help show help for choices
	-a {A,B,C} optional choice (default A)

	Arguments:
	b mandatory choice (choose from {A,B,C})
	`
	msg := builder.String()

	space := regexp.MustCompile(`\s+`)
	msg = space.ReplaceAllString(msg, " ")
	expected = space.ReplaceAllString(expected, " ")

	if msg != expected {
		t.Errorf("unexpected: %s", msg)
		t.Errorf("unexpected: %s", expected)
	}
}

func TestSet(t *testing.T) {
	err := cmd.RunError([]string{"-a", "B", "C"})
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if args.A.Choice != "B" {
		t.Error("A wasn't set correctly")
	}

	if args.B.Choice != "C" {
		t.Error("B wasn't set correctly")
	}

	fmt.Println("next parse...")

	err = cmd.RunError([]string{"-a", "X"})
	switch err {
	case nil:
		t.Error("expected an error")
	default:
		if err.Error() != "parsing flag -a: X is not a valid choice, must be in {A,B,C}" {
			t.Errorf("unexpected error: %s", err)
		}
	}

	builder := new(strings.Builder)
	cmd.SetOutput(builder)
	cmd.SetErrorFlag(failure.ContinueOnError)
	fmt.Println("next parse...")
	cmd.Run([]string{"-a", "X"})

	expected := "Error: parsing flag -a: X is not a valid choice, must be in {A,B,C}.\n"
	msg := builder.String()

	space := regexp.MustCompile(`\s+`)
	msg = space.ReplaceAllString(msg, " ")
	expected = space.ReplaceAllString(expected, " ")

	if msg != expected {
		t.Errorf("unexpected: %s", msg)
		t.Errorf("unexpected: %s", expected)
	}
}
