package abad

import (
	"fmt"

	"github.com/NeowayLabs/abad/ast"
	"github.com/NeowayLabs/abad/parser"
	"github.com/NeowayLabs/abad/types"
	"github.com/NeowayLabs/abad/token"
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

	return a.eval(program)
}

func (a *Abad) eval(n ast.Node) (Obj, error) {
	switch n.Type() {
	case ast.NodeProgram:
		return a.evalProgram(n.(*ast.Program))
	case ast.NodeNumber:
		val := n.(ast.Number)
		return types.Number(val.Value()), nil
	case ast.NodeUnaryExpr:
		expr := n.(*ast.UnaryExpr)
		return a.evalUnaryExpr(expr)
	}

	panic(fmt.Sprintf("AST(%s) not implemented", n))
	return nil, nil
}

func (a *Abad) evalProgram(stmts *ast.Program) (Obj, error) {
	var (
		result Obj
		err    error
	)
	for _, node := range stmts.Nodes {
		result, err = a.eval(node)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

func (a *Abad) evalUnaryExpr(expr *ast.UnaryExpr) (Obj, error) {
	op := expr.Operator
	obj, err := a.eval(expr.Operand)
	if err != nil {
		return nil, err
	}

	// TODO(i4k): UnaryExpr could work in any expression in js
	// examples below are valid:
	//   -[]
	//	 +[]
	//   -{}
	//   +{}
	//   -eval("0")
	//   -new Object()
	num, ok := obj.(types.Number)
	if !ok {
		return nil, fmt.Errorf("not a number: %s", obj)
	}

	switch op {
	case token.Minus:
		num = -num
	case token.Plus:
		num = +num
	default:
		return nil, fmt.Errorf("unsupported unary operator: %s", op)
	}

	return num, nil
}