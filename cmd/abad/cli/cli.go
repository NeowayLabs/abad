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

func NewCli(fname string, in io.Reader, out io.Writer) (*Cli, error) {
	ecma, err := abad.NewAbad(fname)
	if err != nil {
		return nil, err
	}

	return NewWithJS(ecma, in, out), nil
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
		c.error(err)
		return
	}

	line = trimnl(line)

	obj, err := c.js.Eval(line)
	if err != nil {
		c.error(err)
		return
	}

	fmt.Fprintf(c.out, "< %s\n", obj.ToString().String())
}

func (c *Cli) Repl() {
	for {
		c.ReadEval()
	}
}

func (c *Cli) error(err error) {
	fmt.Fprintf(c.out, "%s\n", err)
}

func trimnl(line string) string {
	return line[:len(line)-1]
}