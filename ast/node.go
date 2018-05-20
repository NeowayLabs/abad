package ast

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/NeowayLabs/abad/token"
)

type (
	// Node type
	NodeType int

	// All node types implement the Node interface
	Node interface {
		Type() NodeType
		String() string
		Equal(other Node) bool
	}

	// Program Abstract Syntax Tree
	Program struct {
		Nodes []Node
	}

	Number float64

	// UnaryExpr is a unary expression (-a, +a, ~a, and so on)
	UnaryExpr struct {
		Operator token.Type
		Operand  Node
	}
)

const (
	NodeProgram NodeType = iota + 1

	exprBegin

	NodeNumber
	NodeUnaryExpr

	exprEnd
)

// console.log(Number.EPSILON);
// 2.220446049250313e-16
var ε = math.Pow(2, -52)

func IsExpr(node Node) bool {
	return node.Type() > exprBegin &&
		node.Type() < exprEnd
}

func (_ *Program) Type() NodeType {
	return NodeProgram
}

func (p *Program) String() string {
	var stmts []string
	for _, stmt := range p.Nodes {
		stmts = append(stmts, stmt.String())
	}
	return strings.Join(stmts, "\n")
}

func (p *Program) Equal(other Node) bool {
	if other.Type() != NodeProgram {
		return false
	}

	o := other.(*Program)
	if len(p.Nodes) != len(o.Nodes) {
		return false
	}

	for i := 0; i < len(p.Nodes); i++ {
		if !p.Nodes[i].Equal(o.Nodes[i]) {
			return false
		}
	}
	return true
}

func NewNumber(a float64) Number {
	return Number(a)
}

func NewIntNumber(a int64) Number {
	return Number(float64(a))
}

func (a Number) Value() float64 { return float64(a) }

func (_ Number) Type() NodeType {
	return NodeNumber
}

// TODO(i4k): Implements correct javascript number to string
// representation.
func (a Number) String() string {
	return strconv.FormatFloat(float64(a), 'f', -1, 64)
}

func (a Number) Equal(other Node) bool {
	if other.Type() != NodeNumber {
		return false
	}

	o := other.(Number)
	return floatEquals(float64(a), float64(o))
}

func NewUnaryExpr(operator token.Type, operand Node) *UnaryExpr {
	return &UnaryExpr{
		Operator: operator,
		Operand:  operand,
	}
}

func (_ *UnaryExpr) Type() NodeType {
	return NodeUnaryExpr
}

func (a *UnaryExpr) String() string {
	return fmt.Sprintf("%s%s", a.Operator, a.Operand)
}

func (a *UnaryExpr) Equal(other Node) bool {
	if other.Type() != a.Type() {
		return false
	}

	o := other.(*UnaryExpr)
	if a.Operator != o.Operator {
		return false
	}

	return a.Operand.Equal(o.Operand)
}

func floatEquals(a, b float64) bool {
	return math.Abs(a-b) < ε && math.Abs(b-a) < ε
}