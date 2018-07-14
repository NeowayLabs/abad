package abad

import (
	"fmt"

	"github.com/NeowayLabs/abad/ast"
	"github.com/NeowayLabs/abad/builtins"
	"github.com/NeowayLabs/abad/internal/utf16"
	"github.com/NeowayLabs/abad/parser"
	"github.com/NeowayLabs/abad/token"
	"github.com/NeowayLabs/abad/types"
)

type (
	// Abad interpreter, a very bad one.
	Abad struct {
		global *types.DataObject
	}
)

var (
	consoleAttr = utf16.S("console")
)

// NewAbad creates a new ecma script evaluator.
func NewAbad() (*Abad, error) {
	a := &Abad{}
	return a, a.setup()
}

// Eval the code when no filename is involved (interactive/repl mode).
func (a *Abad) Eval(code string) (types.Value, error) {
	return a.EvalFile("<interactive>", code)
}

// EvalFile the code that was obtained from filename.
func (a *Abad) EvalFile(filename string, code string) (types.Value, error) {
	program, err := parser.Parse(filename, code)
	if err != nil {
		return nil, fmt.Errorf("parser error: %s", err)
	}
	return a.eval(program)
}

func (a *Abad) eval(n ast.Node) (types.Value, error) {
	if ast.IsExpr(n) {
		return a.evalExpr(n)
	}

	var ret types.Value
	var err error

	switch n.Type() {
	case ast.NodeProgram:
		ret, err = a.evalProgram(n.(*ast.Program))
	default:
		panic(fmt.Sprintf("AST(%s) not implemented", n))
	}

	return ret, err
}

func (a *Abad) setup() error {
	console, err := builtins.NewConsole()
	if err != nil {
		return err
	}

	global := types.NewBaseDataObject()
	err = global.Put(consoleAttr, console, true)
	if err != nil {
		return err
	}

	a.global = global
	return nil
}

func (a *Abad) evalProgram(stmts *ast.Program) (types.Value, error) {
	var (
		result types.Value
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

func (a *Abad) evalUnaryExpr(expr *ast.UnaryExpr) (types.Value, error) {
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

func (a *Abad) evalExpr(n ast.Node) (types.Value, error) {
	if !ast.IsExpr(n) {
		panic("internal error: not an expression")
	}

	var ret types.Value
	var err error

	switch n.Type() {
	case ast.NodeNumber:
		val := n.(ast.Number)
		ret, err = types.Number(val.Value()), nil
	case ast.NodeString:	
		val := n.(ast.String)
		ret, err = types.String(val), nil
	case ast.NodeIdent:
		val := n.(ast.Ident)
		ret, err = a.evalIdentExpr(val)
	case ast.NodeMemberExpr:
		val := n.(*ast.MemberExpr)
		ret, err = a.evalMemberExpr(val)
	case ast.NodeCallExpr:
		val := n.(*ast.CallExpr)
		ret, err = a.evalCallExpr(val)
	case ast.NodeUnaryExpr:
		expr := n.(*ast.UnaryExpr)
		ret, err = a.evalUnaryExpr(expr)
	default:
		panic(fmt.Sprintf("unknown node type: %v", n))
	}

	return ret, err
}

func (a *Abad) evalIdentExpr(ident ast.Ident) (types.Value, error) {
	val, err := a.global.Get(utf16.Str(ident))
	if err != nil {
		return nil, err
	}

	if types.StrictEqual(val, types.Undefined) {
		return nil, fmt.Errorf("[%s] is not defined",
			ident.String())
	}

	return val, nil
}

func (a *Abad) evalMemberExpr(member *ast.MemberExpr) (types.Value, error) {
	objval, err := a.evalExpr(member.Object)
	if err != nil {
		return nil, err
	}

	if objval.Kind() != types.KindObject {
		panic("wrapping primitive values not implemented yet")
	}

	obj, err := objval.ToObject()
	if err != nil {
		return nil, err
	}

	return obj.Get(utf16.Str(member.Property))
}

func (a *Abad) evalCallExpr(call *ast.CallExpr) (types.Value, error) {
	// TODO(i4k): safe to assume the AST is ok?
	objval, err := a.evalExpr(call.Callee)
	if err != nil {
		return nil, err
	}

	obj, err := objval.ToObject() // wraps primitives (if needed)
	if err != nil {
		return nil, err
	}

	fun, ok := obj.(types.Function)
	if !ok {
		return nil, fmt.Errorf("%s is not a function", objval.Kind())
	}

	args, err := a.evalArgs(call.Args)
	if err != nil {
		return nil, err
	}

	return fun.Call(obj, args), nil
}

func (a *Abad) evalArgs(args []ast.Node) ([]types.Value, error) {
	var vargs []types.Value

	for _, arg := range args {
		v, err := a.evalExpr(arg)
		if err != nil {
			return nil, err
		}

		vargs = append(vargs, v)
	}

	return vargs, nil
}
