package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/mailund/cli"
)

type args struct {
	In  cli.InFile  `flag:"in" short:"i" descr:"input file"`
	Out cli.OutFile `flag:"out" short:"o" descr:"output file"`
}

func initCat() interface{} {
	return &args{
		In:  cli.InFile{Reader: os.Stdin},
		Out: cli.OutFile{Writer: os.Stdout},
	}
}

func catAction(i interface{}) {
	a, _ := i.(*args)

	buf, err := ioutil.ReadAll(a.In)
	if err != nil {
		log.Fatalf("error reading file: %s", err)
	}

	err = a.In.Close()
	if err != nil {
		log.Fatalf("error closing file: %s", err)
	}

	_, err = a.Out.Write(buf)
	if err != nil {
		log.Fatalf("error writing file: %s", err)
	}

	err = a.Out.Close()
	if err != nil {
		log.Fatalf("error closing file: %s", err)
	}
}

var catCmd = cli.NewCommand(
	cli.CommandSpec{
		Name:   "cat",
		Long:   "Writes the content of one file to another.",
		Init:   initCat,
		Action: catAction,
	})

func main() {
	catCmd.Run(os.Args[1:])
}
