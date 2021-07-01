package cli

import (
	"io"
	"os"

	"github.com/mailund/cli/interfaces"
)

// OutFile represents an open file as an argument
type OutFile struct {
	io.Writer
	fname string
}

// Set implements the PosValue/FlagValue interface
func (o *OutFile) Set(x string) error {
	f, err := os.Create(x)
	if err != nil {
		return interfaces.ParseErrorf("couldn't open file %s: %s", x, err)
	}

	o.Writer = f

	return nil
}

// String implements the FlagValue interface
func (o *OutFile) String() string {
	switch o.Writer {
	case os.Stdout:
		return "stdout"
	case os.Stderr:
		return "stderr"
	default:
		return o.fname
	}
}

// Validate implements the validaator interface. For flags, we either have to
// have a writer or a valid file name.
func (o *OutFile) Validate(flag bool) error {
	if !flag || o.Writer != nil || o.fname != "" {
		return nil // we have a valid default, or we will get an argument
	}

	return interfaces.SpecErrorf("outfile does not have a valid default")
}

// Close implements the io.Closer interface by forwarding to the writer
func (o *OutFile) Close() error {
	if closer, ok := o.Writer.(io.Closer); ok {
		return closer.Close()
	}

	return nil
}
