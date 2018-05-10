package parser_test

import (
	"fmt"
	"testing"

	"github.com/NeowayLabs/abad/ast"
	"github.com/NeowayLabs/abad/parser"
	"github.com/NeowayLabs/abad/token"
	"github.com/madlambda/spells/assert"
)

type testcase struct {
	input       string
	expected    ast.Node
	expectedErr error
}

var E = fmt.Errorf

func TestParserNumbers(t *testing.T) {
	for _, tc := range []testcase{
		{
			input:    "1",
			expected: ast.NewIntNumber(1),
		},
		{
			input:    "1234567890",
			expected: ast.NewIntNumber(1234567890),
		},
		{
			input:       "1a",
			expectedErr: E("tests.js:1:0: invalid token: 1a"),
		},
		{
			input:    "0x0",
			expected: ast.NewIntNumber(0),
		},
		{
			input:    "0x1234567890abcdef",
			expected: ast.NewIntNumber(0x1234567890abcdef),
		},
		{
			input:    "0xff",
			expected: ast.NewIntNumber(0xff),
		},
		{
			input:    ".1",
			expected: ast.NewNumber(0.1),
		},
		{
			input:    ".0000",
			expected: ast.NewNumber(0.0),
		},
		{
			input:    "1234",
			expected: ast.NewIntNumber(1234),
		},
		{
			input:    "0.12345",
			expected: ast.NewNumber(0.12345),
		},
		{
			input:       "0.a",
			expectedErr: E("tests.js:1:0: invalid token: 0.a"),
		},
		{
			input:       "12.13.",
			expectedErr: E("tests.js:1:0: invalid token: 12.13."),
		},
		{
			input:    "1.0e10",
			expected: ast.NewNumber(1.0e10),
		},
		{
			input:    "1e10",
			expected: ast.NewNumber(1e10),
		},
		{
			input:    ".1e10",
			expected: ast.NewNumber(.1e10),
		},
		{
			input:    "1e-10",
			expected: ast.NewNumber(1e-10),
		},
		{
			input: "-1",
			expected: ast.NewUnaryExpr(
				token.Minus, ast.NewNumber(1),
			),
		},
		{
			input: "-1234",
			expected: ast.NewUnaryExpr(
				token.Minus, ast.NewNumber(1234),
			),
		},
		{
			input: "-0x0",
			expected: ast.NewUnaryExpr(
				token.Minus, ast.NewNumber(0),
			),
		},
		{
			input: "-0xff",
			expected: ast.NewUnaryExpr(
				token.Minus, ast.NewNumber(255),
			),
		},
		{
			input:    "-.0",
			expected: ast.NewUnaryExpr(
				token.Minus, ast.NewNumber(0),
			),
		},
		{
			input:    "-.0e1",
			expected: ast.NewUnaryExpr(
				token.Minus, ast.NewNumber(0),
			),
		},
		{
			input:       "-12.13.",
			expectedErr: E("tests.js:1:0: invalid token: 12.13."),
		},
		{
			input:    "-1e-10",
			expected: ast.NewUnaryExpr(
				token.Minus, ast.NewNumber(1.0e-10),
			),
		},
	} {
		tree, err := parser.Parse("tests.js", tc.input)
		assert.EqualErrs(t, tc.expectedErr, err, "parser err")

		if err != nil {
			continue
		}

		nodes := tree.Nodes
		if len(nodes) != 1 {
			t.Fatalf("number tests must be isolated: %v", nodes)
		}

		got := nodes[0]
		if got.Type() != tc.expected.Type() {
			t.Fatalf("literals type differ: %d != %d (%s)",
				got.Type(), tc.expected.Type(), tc.input)
		}

		if !tc.expected.Equal(got) {
			t.Fatalf("Numbers differ: '%s' != '%s'",
				got, tc.expected)
		}
	}
}