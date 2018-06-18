package parser_test

import (
	"fmt"
	"testing"

	"github.com/NeowayLabs/abad/ast"
	"github.com/NeowayLabs/abad/internal/utf16"
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
			input: "-.0",
			expected: ast.NewUnaryExpr(
				token.Minus, ast.NewNumber(0),
			),
		},
		{
			input: "-.0e1",
			expected: ast.NewUnaryExpr(
				token.Minus, ast.NewNumber(0),
			),
		},
		{
			input:       "-12.13.",
			expectedErr: E("tests.js:1:0: invalid token: 12.13."),
		},
		{
			input: "-1e-10",
			expected: ast.NewUnaryExpr(
				token.Minus, ast.NewNumber(1.0e-10),
			),
		},
		{
			input: "-+0",
			expected: ast.NewUnaryExpr(
				token.Minus, ast.NewUnaryExpr(
					token.Plus, ast.NewNumber(0),
				),
			),
		},
		{
			input: "-+-+0",
			expected: ast.NewUnaryExpr(token.Minus,
				ast.NewUnaryExpr(token.Plus,
					ast.NewUnaryExpr(token.Minus,
						ast.NewUnaryExpr(token.Plus,
							ast.NewNumber(0))))),
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

func TestIdentifier(t *testing.T) {
	for _, tc := range []struct {
		input       string
		expected    ast.Node
		expectedErr error
	}{
		{
			input:    "_",
			expected: ast.NewIdent(utf16.S("_")),
		},
		{
			input:    "$",
			expected: ast.NewIdent(utf16.S("$")),
		},
		{
			input:    "console",
			expected: ast.NewIdent(utf16.S("console")),
		},
		{
			input:    "angular",
			expected: ast.NewIdent(utf16.S("angular")),
		},
		{
			input:    "___hyped___",
			expected: ast.NewIdent(utf16.S("___hyped___")),
		},
		{
			input:    "a$b$c",
			expected: ast.NewIdent(utf16.S("a$b$c")),
		},
	} {
		testParser(t, tc.input, tc.expected, tc.expectedErr)
	}
}

func TestMemberExpr(t *testing.T) {
	for _, tc := range []struct {
		input       string
		expected    ast.Node
		expectedErr error
	}{
		{
			input: "console.log",
			expected: ast.NewMemberExpr(
				ast.NewIdent(utf16.S("console")),
				ast.NewIdent(utf16.S("log")),
			),
		},
		{
			input:       "console.",
			expectedErr: E("tests.js:1:0: unexpected EOF"),
		},
		{
			input: "self.a",
			expected: ast.NewMemberExpr(
				ast.NewIdent(utf16.S("self")),
				ast.NewIdent(utf16.S("a")),
			),
		},
		{
			input: "self.self.self", // same as: (self.self).self
			expected: ast.NewMemberExpr(
				ast.NewMemberExpr(ast.NewIdent(utf16.S("self")), ast.NewIdent(utf16.S("self"))),
				ast.NewIdent(utf16.S("self")),
			),
		},
		{
			input: "a.b.c.d.e.f", // same as: ((((a.b).c).d).e).f)
			expected: ast.NewMemberExpr(
				ast.NewMemberExpr(
					ast.NewMemberExpr(
						ast.NewMemberExpr(
							ast.NewMemberExpr(ast.NewIdent(utf16.S("a")), ast.NewIdent(utf16.S("b"))),
							ast.NewIdent(utf16.S("c")),
						),
						ast.NewIdent(utf16.S("d")),
					),
					ast.NewIdent(utf16.S("e")),
				),
				ast.NewIdent(utf16.S("f")),
			),
		},
	} {
		testParser(t, tc.input, tc.expected, tc.expectedErr)
	}
}

func TestParserFuncall(t *testing.T) {
	for _, tc := range []struct {
		input       string
		expected    ast.Node
		expectedErr error
	}{
		{
			input: "console.log()",
			expected: ast.NewCallExpr(
				ast.NewMemberExpr(
					ast.NewIdent(utf16.S("console")),
					ast.NewIdent(utf16.S("log")),
				),
				[]ast.Node{},
			),
		},
		{
			input: "console.log(2.0)",
			expected: ast.NewCallExpr(
				ast.NewMemberExpr(
					ast.NewIdent(utf16.S("console")),
					ast.NewIdent(utf16.S("log")),
				),
				[]ast.Node{ast.NewNumber(2.0)},
			),
		},
		{
			input: "self.console.log(2.0)",
			expected: ast.NewCallExpr(
				ast.NewMemberExpr(
					ast.NewMemberExpr(
						ast.NewIdent(utf16.S("self")),
						ast.NewIdent(utf16.S("console")),
					),
					ast.NewIdent(utf16.S("log")),
				),
				[]ast.Node{ast.NewNumber(2.0)},
			),
		},
	} {
		testParser(t, tc.input, tc.expected, tc.expectedErr)
	}
}

func testParser(
	t *testing.T, input string, expected ast.Node, expectedErr error,
) {
	tree, err := parser.Parse("tests.js", input)
	assert.EqualErrs(t, expectedErr, err, "parser err")

	if err != nil {
		return
	}

	nodes := tree.Nodes
	if len(nodes) != 1 {
		t.Fatalf("memberexpr tests must be isolated: %v", nodes)
	}

	got := nodes[0]
	if got.Type() != expected.Type() {
		t.Fatalf("type differ: %d != %d (%s)",
			got.Type(), expected.Type(), input)
	}

	if !expected.Equal(got) {
		t.Fatalf("Identifier differ: '%s' != '%s'",
			got, expected)
	}
}
