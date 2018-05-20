package main

import (
	"fmt"
	"os"

	"github.com/NeowayLabs/abad/cmd/abad/cli"
)

const filename = "<anonymous>"

func run() error {
	cli, err := cli.NewCli(filename, os.Stdin, os.Stdout)
	if err != nil {
		return err
	}

	cli.Repl()
	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Printf("error: %s\n", err)
		os.Exit(1)
	}
}