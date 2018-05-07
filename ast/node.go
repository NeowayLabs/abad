package ast

import (
	"math"
	"strconv"
	"strings"
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
)

const (
	NodeProgram NodeType = iota + 1
	NodeNumber
)

// copied from V8 running:
// console.log(Number.EPSILON);
// TODO(i4k): Inspect v8 source code for the right value.
const ε = 2.220446049250313e-16

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

func floatEquals(a, b float64) bool {
	return (math.Abs(a-b) < ε && math.Abs(b-a) < ε)
}