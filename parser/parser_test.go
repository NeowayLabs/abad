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

var E = fmt.Errorf

func TestParserNumbers(t *testing.T) {

	runTests(t, []TestCase{
		{
			name: "SmallDecimal",
			code: "1",
			want: ast.NewIntNumber(1),
		},
		{
			name: "BigDecimal",
			code: "1234567890",
			want: ast.NewIntNumber(1234567890),
		},
		{
			name:    "InvalidDecimal",
			code:    "1a",
			wantErr: E("tests.js:1:0: invalid token: 1a"),
		},
		{
			name: "SmallHexadecimal",
			code: "0x0",
			want: ast.NewIntNumber(0),
		},
		{
			name: "BigHexaDecimal",
			code: "0x1234567890abcdef",
			want: ast.NewIntNumber(0x1234567890abcdef),
		},
		{
			name: "HexadecimalFF",
			code: "0xff",
			want: ast.NewIntNumber(0xff),
		},
		{
			name: "SmallRealNumber",
			code: ".1",
			want: ast.NewNumber(0.1),
		},
		{
			name: "ZeroRealNumer",
			code: ".0000",
			want: ast.NewNumber(0.0),
		},
		{
			name: "SomeDecimal",
			code: "1234",
			want: ast.NewIntNumber(1234),
		},
		{
			name: "SmallRealNumberWithMultipleDigits",
			code: "0.12345",
			want: ast.NewNumber(0.12345),
		},
		{
			name:    "InvalidRealNumberWithLetter",
			code:    "0.a",
			wantErr: E("tests.js:1:0: invalid token: 0.a"),
		},
		{
			name:    "InvalidRealNumberWithTwoDots",
			code:    "12.13.",
			wantErr: E("tests.js:1:0: invalid token: 12.13."),
		},
		{
			name: "RealNumberWithExponent",
			code: "1.0e10",
			want: ast.NewNumber(1.0e10),
		},
		{
			name: "DecimalWithExponent",
			code: "1e10",
			want: ast.NewNumber(1e10),
		},
		{
			name: "SmallRealNumberWithExponent",
			code: ".1e10",
			want: ast.NewNumber(.1e10),
		},
		{
			name: "DecimalWithNegativeExponent",
			code: "1e-10",
			want: ast.NewNumber(1e-10),
		},
		{
			name: "NegativeDecimalWithOneDigit",
			code: "-1",
			want: ast.NewUnaryExpr(
				token.Minus, ast.NewNumber(1),
			),
		},
		{
			name: "NegativeDecimalWithMultipleDigits",
			code: "-1234",
			want: ast.NewUnaryExpr(
				token.Minus, ast.NewNumber(1234),
			),
		},
		{
			name: "NegativeZeroHexadecimal",
			code: "-0x0",
			want: ast.NewUnaryExpr(
				token.Minus, ast.NewNumber(0),
			),
		},
		{
			name: "NegativeFFHexadecimal",
			code: "-0xff",
			want: ast.NewUnaryExpr(
				token.Minus, ast.NewNumber(255),
			),
		},
		{
			name: "NegativeZeroRealNumber",
			code: "-.0",
			want: ast.NewUnaryExpr(
				token.Minus, ast.NewNumber(0),
			),
		},
		{
			name: "NegativeZeroRealNumberWithExponent",
			code: "-.0e1",
			want: ast.NewUnaryExpr(
				token.Minus, ast.NewNumber(0),
			),
		},
		{
			name:    "InvalidNegativeRealNumber",
			code:    "-12.13.",
			wantErr: E("tests.js:1:0: invalid token: 12.13."),
		},
		{
			name: "NegativeDecimalWithNegativeExponent",
			code: "-1e-10",
			want: ast.NewUnaryExpr(
				token.Minus, ast.NewNumber(1.0e-10),
			),
		},
		{
			name: "NegativePlusZeroDecimal",
			code: "-+0",
			want: ast.NewUnaryExpr(
				token.Minus, ast.NewUnaryExpr(
					token.Plus, ast.NewNumber(0),
				),
			),
		},
		{
			name: "InterleavedNegativeWithPlusAndZeroDecimal",
			code: "-+-+0",
			want: ast.NewUnaryExpr(token.Minus,
				ast.NewUnaryExpr(token.Plus,
					ast.NewUnaryExpr(token.Minus,
						ast.NewUnaryExpr(token.Plus,
							ast.NewNumber(0))))),
		},
	})
}

func TestIdentifier(t *testing.T) {
	runTests(t, []TestCase{
		{
			name: "Underscore",
			code: "_",
			want: ast.NewIdent(utf16.S("_")),
		},
		{
			name: "Dolar",
			code: "$",
			want: ast.NewIdent(utf16.S("$")),
		},
		{
			name: "Console",
			code: "console",
			want: ast.NewIdent(utf16.S("console")),
		},
		{
			name: "AngularSux",
			code: "angular",
			want: ast.NewIdent(utf16.S("angular")),
		},
		{
			name: "HyperdUnderscores",
			code: "___hyped___",
			want: ast.NewIdent(utf16.S("___hyped___")),
		},
		{
			name: "LettersAndDolars",
			code: "a$b$c",
			want: ast.NewIdent(utf16.S("a$b$c")),
		},
	})
}

func TestMemberExpr(t *testing.T) {
	runTests(t, []TestCase{
		{
			name: "AccessingLogOnConsole",
			code: "console.log",
			want: ast.NewMemberExpr(
				ast.NewIdent(utf16.S("console")),
				ast.NewIdent(utf16.S("log")),
			),
		},
		{
			name:    "ErrorAccessingEmptyMember",
			code:    "console.",
			wantErr: E("tests.js:1:0: unexpected EOF"),
		},
		{
			name: "AccessMemberOfSelf",
			code: "self.a",
			want: ast.NewMemberExpr(
				ast.NewIdent(utf16.S("self")),
				ast.NewIdent(utf16.S("a")),
			),
		},
		{
			name: "OneLevelOfNesting",
			code: "self.self.self", // same as: (self.self).self
			want: ast.NewMemberExpr(
				ast.NewMemberExpr(ast.NewIdent(utf16.S("self")), ast.NewIdent(utf16.S("self"))),
				ast.NewIdent(utf16.S("self")),
			),
		},
		{
			name: "MultipleLevelsOfNesting",
			code: "a.b.c.d.e.f", // same as: ((((a.b).c).d).e).f)
			want: ast.NewMemberExpr(
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
	})
}

func TestParserFuncall(t *testing.T) {

	runTests(t, []TestCase{
		{
			name: "OneParameter",
			code: "a()",
			want: ast.NewCallExpr(
				ast.NewIdent(utf16.S("a")),
				[]ast.Node{},
			),
		},
		{
			name: "MemberAccessWithoutParams",
			code: "console.log()",
			want: ast.NewCallExpr(
				ast.NewMemberExpr(
					ast.NewIdent(utf16.S("console")),
					ast.NewIdent(utf16.S("log")),
				),
				[]ast.Node{},
			),
		},
		//{
		//	name: "MultipleCallsNoParams",
		//	code: "a();\nb();",
		//	wants: []ast.Node{
		//		ast.NewCallExpr(
		//			ast.NewIdent(utf16.S("a")),
		//			[]ast.Node{},
		//		),
		//		ast.NewCallExpr(
		//			ast.NewIdent(utf16.S("b")),
		//			[]ast.Node{},
		//		),
		//	},
		//},
		{
			name: "MemberAccessWithDecimalParam",
			code: "console.log(2.0)",
			want: ast.NewCallExpr(
				ast.NewMemberExpr(
					ast.NewIdent(utf16.S("console")),
					ast.NewIdent(utf16.S("log")),
				),
				[]ast.Node{ast.NewNumber(2.0)},
			),
		},
		{
			name: "NestedMemberAccessWithDecimalParam",
			code: "self.console.log(2.0)",
			want: ast.NewCallExpr(
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
	})
}

// TestCase is the description of an parser related test.
// The fields want and wants are mutually exclusive, you should
// never provide both. If "wants" is provided the "want" field will be ignored.
//
// This is supposed to make it easier to test single nodes and multiple nodes.
type TestCase struct {
	name    string
	code    string
	want    ast.Node
	wants   []ast.Node
	wantErr error
}

func (tc *TestCase) run(t *testing.T) {
	t.Run(tc.name, func(t *testing.T) {
		tree, err := parser.Parse("tests.js", tc.code)
		assert.EqualErrs(t, tc.wantErr, err, "parser err")

		if err != nil {
			return
		}

		if tc.wants == nil {
			assertEqualNodes(t, []ast.Node{tc.want}, tree.Nodes)
			return
		}

		assertEqualNodes(t, tc.wants, tree.Nodes)
	})
}

func runTests(t *testing.T, tcases []TestCase) {
	for _, tcase := range tcases {
		tcase.run(t)
	}
}

func assertEqualNodes(t *testing.T, want []ast.Node, got []ast.Node) {
	if len(want) != len(got) {
		t.Errorf("want[%d] nodes but got[%d] nodes", len(want), len(got))
		t.Fatalf("want[%v] != got[%v]", want, got)
	}

	for i, w := range want {
		g := got[i]
		if !w.Equal(g) {
			t.Errorf("wanted node[%d][%v] != got node[%d][%v]", i, w, i, g)
		}
	}
}
