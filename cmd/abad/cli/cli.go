package cli

import (
	"bufio"
	"fmt"
	"io"

	"github.com/NeowayLabs/abad"
)

type (
	Cli struct {
		in  io.Reader
		out io.Writer

		js *abad.Abad
	}
)

func NewCli(fname string, in io.Reader, out io.Writer) *Cli {
	return NewWithJS(abad.NewAbad(fname), in, out)
}

func NewWithJS(js *abad.Abad, in io.Reader, out io.Writer) *Cli {
	return &Cli{
		in:  in,
		out: out,
		js:  js,
	}
}

func (c *Cli) ReadEval() {
	fmt.Fprintf(c.out, "> ")
	bio := bufio.NewReader(c.in)
	line, err := bio.ReadString('\n')
	if err != nil {
		c.errorf(err)
		return
	}

	if len(line) > 0 {
		line = line[:len(line)-1]
	}

	obj, err := c.js.Eval(line)
	if err != nil {
		c.errorf(err)
		return
	}

	fmt.Fprintf(c.out, "< %s\n", obj)
}

func (c *Cli) Repl() {
	for {
		c.ReadEval()
	}
}

func (c *Cli) errorf(err error) {
	fmt.Fprintf(c.out, "%s\n", err)
}