package abad

import (
	"fmt"

	"github.com/NeowayLabs/abad/ast"
	"github.com/NeowayLabs/abad/parser"
	"github.com/NeowayLabs/abad/types"
)

type (
	// Obj types
	Obj interface {
		String() string
	}

	// Abad interpreter, a very bad one.
	Abad struct {
		filename string
	}
)

func NewAbad(filename string) *Abad {
	return &Abad{
		filename: filename,
	}
}

func (a *Abad) Eval(code string) (Obj, error) {
	program, err := parser.Parse(a.filename, code)
	if err != nil {
		return nil, err
	}

	var result Obj

	for _, node := range program.Nodes {
		switch node.Type() {
		case ast.NodeNumber:
			val := node.(ast.Number)
			result = types.Number(val.Value())
		default:
			panic(fmt.Sprintf("AST(%s) not implemented", node))
		}
	}

	return result, nil
}