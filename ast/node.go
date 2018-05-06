package ast

type (
	// Node type
	Type int

	// All node types implement the Node interface
	Node interface {
		Type() Type
		String() string
	}

	// Abstract Syntax Tree
	AST struct {
		Nodes []Node
	}
)