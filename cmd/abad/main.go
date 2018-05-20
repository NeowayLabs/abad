package main

import (
	"os"

	"github.com/NeowayLabs/abad/cmd/abad/cli"
)

const filename = "<anonymous>"

func run() {
	cli := cli.NewCli(filename, os.Stdin, os.Stdout)
	cli.Repl()
}

func main() {
	run()
}