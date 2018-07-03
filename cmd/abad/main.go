package main

import (
	"io/ioutil"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/NeowayLabs/abad"
	"github.com/NeowayLabs/abad/cmd/abad/cli"
)


const defaultFilename = "<anonymous>"


func repl() error {

	cli, err := cli.NewCli(defaultFilename, os.Stdin, os.Stdout)
	if err != nil {
		return err
	}

	cli.Repl()
	return nil
}

func eval(codepath string) error {
	code, err := ioutil.ReadFile(codepath)
	if err != nil {
		return err
	}
	return evalCode(filepath.Base(codepath), string(code))
}

func evalCode(filename string, code string) error {
	abadjs, err := abad.NewAbad(filename)
	if err != nil {
		return err
	}
	_, err = abadjs.Eval(code)
	return err
}

func main() {
	var execute string
	var help bool
	
	flag.BoolVar(&help, "help", false, "prints usage")
	flag.StringVar(&execute, "e", "", "execute code")
	flag.Parse()
	
	if help {
		fmt.Println("Abad: the bad JS interpreter")
		flag.PrintDefaults()
		return
	}
	
	if execute != "" {
		abortonerr(evalCode(defaultFilename, execute))
		return
	}
	
	if len(flag.Args()) == 0 {
		abortonerr(repl())
		return
	}
	
	filepath := flag.Args()[0]
	abortonerr(eval(filepath))
}

func abortonerr(err error) {
	if err != nil {
		fmt.Printf("error: %s\n", err)
		os.Exit(1)
	}
}