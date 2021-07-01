package cli_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/mailund/cli"
	"github.com/mailund/cli/interfaces"
	"github.com/mailund/cli/internal/failure"
)

type prepareSuccess struct{}

func (p *prepareSuccess) Set(string) error    { return nil }
func (p *prepareSuccess) String() string      { return "" }
func (p *prepareSuccess) PrepareValue() error { return nil }

type prepareFailure struct{}

func (p *prepareFailure) Set(string) error { return nil }
func (p *prepareFailure) String() string   { return "" }
func (p *prepareFailure) PrepareValue() error {
	return interfaces.ParseErrorf("fail")
}

type varPrepareSuccess struct{}

func (p *varPrepareSuccess) Set([]string) error  { return nil }
func (p *varPrepareSuccess) String() string      { return "" }
func (p *varPrepareSuccess) PrepareValue() error { return nil }

type success struct {
	A prepareSuccess    `flag:"flag" short:"f"`
	B prepareSuccess    `pos:"pos"`
	C varPrepareSuccess `pos:"var"`
}

type flagFail struct {
	A prepareFailure `flag:"flag" short:"f"`
}

func TestSuccess(t *testing.T) {
	builder := new(strings.Builder)
	cmd, err := cli.NewCommandError(
		cli.CommandSpec{
			Init:   func() interface{} { return new(success) },
			Action: func(interface{}) { fmt.Fprintln(builder, "success") },
		})

	if err != nil {
		t.Fatalf("Unexpected construction error: %s", err)
	}

	cmd.SetOutput(builder)
	cmd.Run([]string{"b", "c"})

	if expected, msg := "success\n", builder.String(); msg != expected {
		t.Errorf("unexpected msg: %s", msg)
	}
}

func TestFlagFail(t *testing.T) {
	builder := new(strings.Builder)
	cmd, err := cli.NewCommandError(
		cli.CommandSpec{
			Init:   func() interface{} { return new(flagFail) },
			Action: func(interface{}) { fmt.Fprintln(builder, "success") },
		})

	if err != nil {
		t.Fatalf("Unexpected construction error: %s", err)
	}

	failed := false
	failure.Failure = func() { failed = true }

	err = cmd.RunError([]string{})
	if err == nil {
		t.Fatal("expected an error")
	}

	if err.Error() != "error in flag -f,--flag: fail" {
		t.Errorf("(1) unexpected err msg: %s", err)
	}

	cmd.SetOutput(builder)
	cmd.Run([]string{})

	if !failed {
		t.Error("Expected a failure")
	}

	if expected, msg := "Error: error in flag -f,--flag: fail.\n", builder.String(); msg != expected {
		t.Errorf("(2) unexpected msg: %s", msg)
	}
}
