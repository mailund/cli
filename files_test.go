package cli_test

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/mailund/cli"
	"github.com/mailund/cli/interfaces"
)

type OutArgs struct {
	Flag cli.OutFile `flag:"flag"`
	Pos  cli.OutFile `pos:"pos"`
}

func outArgsCmd(flag, pos cli.OutFile) (*cli.Command, error) {
	return cli.NewCommandError(
		cli.CommandSpec{
			Init: func() interface{} { return &OutArgs{Flag: flag, Pos: pos} },
		})
}

func TestOutFile(t *testing.T) {
	var args = []struct {
		name      string
		flag, pos cli.OutFile
		err       error
	}{
		{
			name: "empty files",
			flag: cli.OutFile{},
			pos:  cli.OutFile{},
			err:  interfaces.SpecErrorf("outfile does not have a valid default"),
		},
		{
			name: "default flag stdout",
			flag: cli.OutFile{Writer: os.Stdout},
			pos:  cli.OutFile{},
			err:  nil,
		},
		{
			name: "default flag stderr", // just to get that part of String()
			flag: cli.OutFile{Writer: os.Stderr},
			pos:  cli.OutFile{},
			err:  nil,
		},
	}

	for _, a := range args {
		t.Run(a.name, func(t *testing.T) {
			_, err := outArgsCmd(a.flag, a.pos)

			switch {
			case err == nil && a.err == nil:
				return
			case err != nil && a.err == nil:
				t.Errorf("Unexpected error: %s", err)

			case err == nil && a.err != nil:
				t.Errorf("Expected an error")

			case err.Error() != a.err.Error():
				t.Errorf("unexpected error, %s vs %s", err, a.err)
			}
		})
	}
}

func TestOutFileSet(t *testing.T) {
	var f cli.OutFile

	// This is a bit of a hack to get a temp file
	file, err := ioutil.TempFile("", "cli")
	if err != nil {
		log.Fatal(err)
	}

	defer os.Remove(file.Name())

	if err = f.Set(file.Name()); err != nil {
		t.Error("We should not get an error opening this file")
	}

	fmt.Fprintln(f, "hello, world!")

	if err = f.Close(); err != nil {
		t.Error("We should not get an error closing this file")
	}

	fcontent, err := ioutil.ReadFile(file.Name())
	if err != nil {
		t.Error("There should be a file to read")
	}

	if string(fcontent) != "hello, world!\n" {
		t.Errorf("File content is wrong: %s", string(fcontent))
	}

	// If we can open the empty file, something is defintely wrong on the system
	if err := f.Set(""); err == nil {
		t.Error("We should *not* be able to open a file with an empty name!")
	}
}

func TestOutFileCLoseStd(t *testing.T) {
	var f = cli.OutFile{Writer: os.Stdout}

	if err := f.Close(); err != nil {
		t.Error("we should be able to 'close' stdout")
	}
}

func TestOutFileCLoseNil(t *testing.T) {
	var f = cli.OutFile{}

	// we should actually also be able to close nil, although
	// I don't think there is ever a situation where we want to
	if err := f.Close(); err != nil {
		t.Error("we should be able to 'close' nil")
	}
}

func TestOutFilePrepare(t *testing.T) {
	var f cli.OutFile

	if err := f.PrepareValue(); err == nil {
		t.Error("this should be an error")
	}

	f.Writer = os.Stdout
	if err := f.PrepareValue(); err != nil {
		t.Error("this should not")
	}

	// This is a bit of a hack to get a temp file
	file, err := ioutil.TempFile("", "cli")
	if err != nil {
		log.Fatal(err)
	}

	defer os.Remove(file.Name())

	f.Writer = nil
	f.Fname = file.Name()

	if err := f.PrepareValue(); err != nil {
		t.Errorf("expected to open %s", f.Fname)
	}

	if f.Writer == nil {
		t.Error("we should have an open file now")
	}

	f.Close()

	// If we can open the empty file, something is defintely wrong on the system
	f.Fname = ""
	if err := f.PrepareValue(); err == nil {
		t.Error("We should *not* be able to open a file with an empty name!")
	}
}

func TestWriteRead(t *testing.T) { //nolint:funlen // it is not too long
	type (
		write struct {
			Out cli.OutFile `flag:"out"`
		}
		read struct {
			In cli.InFile `flag:"in"`
		}
	)

	file, err := ioutil.TempFile("", "cli")
	if err != nil {
		log.Fatal(err)
	}

	defer os.Remove(file.Name())
	tmpname := file.Name()
	res := ""
	writeCmd := cli.NewCommand(
		cli.CommandSpec{
			Name: "write",
			Init: func() interface{} {
				return &write{Out: cli.OutFile{Fname: tmpname}}
			},
			Action: func(i interface{}) {
				out := i.(*write).Out
				if out.Writer == nil {
					t.Fatal("we should have an open file now")
				}
				fmt.Fprintln(out, "hello, world!")
				out.Close()
			},
		})
	readCmd := cli.NewCommand(
		cli.CommandSpec{
			Name: "read",
			Init: func() interface{} {
				return &read{In: cli.InFile{Fname: tmpname}}
			},
			Action: func(i interface{}) {
				in := i.(*read).In
				buf, err := ioutil.ReadAll(in) //nolint:govet // I know err is shadowed
				if err != nil {
					t.Fatalf("error reading file: %s", err)
				}
				fmt.Printf("'%s'\n", string(buf))
				res = string(buf)
			},
		})

	cmd := cli.NewMenu("cmd", "", "", writeCmd, readCmd)

	// write hello world
	err = cmd.RunError([]string{"write"})
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	// now call reader to get the result back
	err = cmd.RunError([]string{"read"})
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if res != "hello, world!\n" {
		t.Errorf("somehow, we got the wrong data: %s", res)
	}
}

func TestInFileUsage(t *testing.T) {
	var f = cli.InFile{Reader: os.Stdin}

	if f.String() != "stdin" {
		t.Error("we should get stdin")
	}

	f.Reader = nil
	f.Fname = "foo"

	if f.String() != `"foo"` {
		t.Errorf("we should get foo but got %s", f.String())
	}
}

func TestInFileCLoseStd(t *testing.T) {
	var f = cli.OutFile{Writer: os.Stdin}

	if err := f.Close(); err != nil {
		t.Error("we should be able to 'close' stdin")
	}
}

func TestInFileCLoseNil(t *testing.T) {
	var f = cli.InFile{}

	// we should actually also be able to close nil, although
	// I don't think there is ever a situation where we want to
	if err := f.Close(); err != nil {
		t.Error("we should be able to 'close' nil")
	}
}

func TestInFilePrepare(t *testing.T) {
	var f cli.InFile

	if err := f.Validate(true); err == nil {
		t.Error("this should be an error")
	}

	f.Reader = os.Stdin
	if err := f.PrepareValue(); err != nil {
		t.Errorf("unexpected error %s", err)
	}

	// This is a bit of a hack to get a temp file
	file, err := ioutil.TempFile("", "cli")
	if err != nil {
		log.Fatal(err)
	}

	defer os.Remove(file.Name())

	if err := f.Set(file.Name()); err != nil {
		t.Errorf("we should be able to open tmpfile: %s", err)
	}

	f.Close()

	f.Fname = file.Name()
	if err := f.PrepareValue(); err != nil {
		t.Errorf("expected to open %s", f.Fname)
	}

	if f.Reader == nil {
		t.Error("we should have an open file now")
	}

	f.Close()

	// If we can open the empty file, something is defintely wrong on the system
	f.Fname = ""
	if err := f.PrepareValue(); err == nil {
		t.Error("We should *not* be able to open a file with an empty name!")
	}
}
