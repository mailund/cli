package main

import (
	"fmt"
	"os"

	"github.com/mailund/cli"
)

func main() {
	cmd := cli.NewCommand(
		cli.CommandSpec{
			Name:   "hello",
			Long:   "Prints hello world",
			Action: func(_ interface{}) { fmt.Println("hello, world!") },
		})

	cmd.Run(os.Args[1:])
}
