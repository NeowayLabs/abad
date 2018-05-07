package parser_test

import (
	"testing"

	"github.com/NeowayLabs/abad/ast"
	"github.com/NeowayLabs/abad/parser"
)

type testcase struct {
	input       string
	expected    ast.Node
	expectedErr string
}

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
			expectedErr: "tests.js:1:0: invalid token: 1a",
		},
		{
			input:    "0x0",
			expected: ast.NewIntNumber(0),
		},
	} {
		tree, err := parser.Parse("tests.js", tc.input)
		if err != nil {
			if tc.expectedErr == "" {
				t.Fatal(err)
			} else if err.Error() != tc.expectedErr {
				t.Fatalf("error differs: Expected [%s] but got [%s]",
					tc.expectedErr, err)
			}

			return
		}

		nodes := tree.Nodes
		if len(nodes) != 1 {
			t.Fatalf("number tests must be isolated: %v", nodes)
		}

		got := nodes[0]
		if got.Type() != tc.expected.Type() {
			t.Fatalf("literals type differ: %d != %d",
				got.Type(), tc.expected.Type())
		}

		if !tc.expected.Equal(got) {
			t.Fatalf("Numbers differ: '%s' != '%s'",
				got, tc.expected)
		}
	}
}