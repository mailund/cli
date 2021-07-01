package cli

import (
	"io"
	"os"

	"github.com/mailund/cli/interfaces"
)

// OutFile represents an open file as an argument
type OutFile struct {
	io.Writer
	Fname string
}

func (o *OutFile) open(fname string) error {
	f, err := os.Create(fname)
	if err != nil {
		return interfaces.ParseErrorf("couldn't open file %s: %s", fname, err)
	}

	o.Writer = f

	return nil
}

// Set implements the PosValue/FlagValue interface
func (o *OutFile) Set(fname string) error {
	return o.open(fname)
}

// String implements the FlagValue interface
func (o *OutFile) String() string {
	switch o.Writer {
	case os.Stdout:
		return "stdout"
	case os.Stderr:
		return "stderr"
	default:
		return `"` + o.Fname + `"`
	}
}

// Validate implements the validaator interface. For flags, we either have to
// have a writer or a valid file name.
func (o *OutFile) Validate(flag bool) error {
	if !flag || o.Writer != nil || o.Fname != "" {
		return nil // we have a valid default, or we will get an argument
	}

	return interfaces.SpecErrorf("outfile does not have a valid default")
}

// Close implements the io.Closer interface by forwarding to the writer
func (o *OutFile) Close() error {
	var err error

	if closer, ok := o.Writer.(io.Closer); ok {
		err = closer.Close()
	}

	o.Writer = nil

	return err
}

// PrepareValue implements the prepare protocol. It will open
// a file if we don't already have a writer but we have a name.
func (o *OutFile) PrepareValue() error {
	if o.Writer != nil {
		return nil // we already have a writer
	}

	return o.open(o.Fname)
}

// InFile represents an open file as an argument
type InFile struct {
	io.Reader
	Fname string
}

func (in *InFile) open(fname string) error {
	f, err := os.Open(fname)
	if err != nil {
		return interfaces.ParseErrorf("couldn't open file %s: %s", fname, err)
	}

	in.Reader = f

	return nil
}

// Set implements the PosValue/FlagValue interface
func (in *InFile) Set(fname string) error {
	return in.open(fname)
}

// String implements the FlagValue interface
func (in *InFile) String() string {
	switch in.Reader {
	case os.Stdin:
		return "stdin"
	default:
		return `"` + in.Fname + `"`
	}
}

// Validate implements the validaator interface. For flags, we either have to
// have a writer or a valid file name.
func (in *InFile) Validate(flag bool) error {
	if !flag || in.Reader != nil || in.Fname != "" {
		return nil // we have a valid default, or we will get an argument
	}

	return interfaces.SpecErrorf("infile does not have a valid default")
}

// Close implements the io.Closer interface by forwarding to the writer
func (in *InFile) Close() error {
	var err error

	if closer, ok := in.Reader.(io.Closer); ok {
		err = closer.Close()
	}

	in.Reader = nil

	return err
}

// PrepareValue implements the prepare protocol. It will open
// a file if we don't already have a writer but we have a name.
func (in *InFile) PrepareValue() error {
	if in.Reader != nil {
		return nil // we already have a reader
	}

	return in.open(in.Fname)
}
