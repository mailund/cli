package cli_test

import (
	"strings"
	"testing"

	"github.com/mailund/cli"
)

type FlagsCallbacks struct {
	A int // not part of the specs

	Fcb  func(string) error              `flag:"cb"`
	Fcbi func(string, interface{}) error `flag:"cbi"`
}

type PosCallbacks1 struct {
	Pcb func(string) error   `pos:"cb"`
	Vcb func([]string) error `pos:"vcb"`
}

type PosCallbacks2 struct {
	A, B int // not part of the specs

	Pcbi func(string, interface{}) error   `pos:"cbi"`
	Vcbi func([]string, interface{}) error `pos:"vcbi"`
}

func TestFlagCallbacks(t *testing.T) {
	var (
		x, y string
	)

	args := &FlagsCallbacks{
		Fcb: func(s string) error {
			x = s
			return nil
		},
		Fcbi: func(s string, i interface{}) error {
			args, ok := i.(*FlagsCallbacks)
			if !ok {
				t.Fatal("Fcbi was not called with the right interface")
			}

			y = s
			args.A = 13

			return nil
		},
	}

	cmd := cli.NewCommand(
		cli.CommandSpec{
			Init: func() interface{} { return args },
		})

	cmd.Run([]string{"--cb=foo", "--cbi=bar"})

	if x != "foo" {
		t.Error("x was not set correctly")
	}

	if y != "bar" {
		t.Error("y was not set correctly")
	}

	if args.A != 13 {
		t.Error("A was not set correctly")
	}
}

func TestPosCallbacks1(t *testing.T) {
	var (
		x, y string
	)

	args := &PosCallbacks1{
		Pcb: func(s string) error {
			x = s
			return nil
		},
		Vcb: func(s []string) error {
			y = strings.Join(s, "")
			return nil
		},
	}

	cmd := cli.NewCommand(
		cli.CommandSpec{
			Init: func() interface{} { return args },
		})

	cmd.Run([]string{"foo", "bar", "baz"})

	if x != "foo" {
		t.Error("x was not set correctly")
	}

	if y != "barbaz" {
		t.Error("y was not set correctly")
	}
}

func TestPosCallbacks2(t *testing.T) {
	var (
		x, y string
	)

	args := &PosCallbacks2{
		Pcbi: func(s string, i interface{}) error {
			args, ok := i.(*PosCallbacks2)
			if !ok {
				t.Fatal("pcbi called with wrong interface")
			}

			args.A = 13
			x = s

			return nil
		},
		Vcbi: func(s []string, i interface{}) error {
			args, ok := i.(*PosCallbacks2)
			if !ok {
				t.Fatal("vcbi called with wrong interface")
			}

			args.B = 42
			y = strings.Join(s, "")

			return nil
		},
	}

	cmd := cli.NewCommand(
		cli.CommandSpec{
			Init: func() interface{} { return args },
		})

	cmd.Run([]string{"foo", "bar", "baz"})

	if x != "foo" {
		t.Error("x was not set correctly")
	}

	if y != "barbaz" {
		t.Error("y was not set correctly")
	}

	if args.A != 13 {
		t.Error("A was not set correctly")
	}

	if args.B != 42 {
		t.Error("B was not set correctly")
	}
}
