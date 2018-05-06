package parser

import (
	"github.com/NeowayLabs/abad/token"
	"github.com/NeowayLabs/abad/lexer"
)

// Parse input source into an AST representation.
func Parse(code string) (ast.AST, error) {
	return parse(utf16.Encode(code))
}

func parse(code utf16.Str) (ast.AST, error) {
	return nil, nil
}