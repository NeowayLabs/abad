package lexer_test

import (
	"testing"
	
	"github.com/NeowayLabs/abad/token"
	"github.com/NeowayLabs/abad/lexer"
	"github.com/NeowayLabs/abad/internal/utf16"
)


type TestCase struct {
	name string
	code utf16.Str
	want []lexer.Tokval
}

var Str func(string) utf16.Str = utf16.S

func TestNumericLiterals(t *testing.T) {

	runTests(t, []TestCase{
		{
			name: "Zero",
			code: Str("0"),
			want: []lexer.Tokval{
				{
					Type: token.Decimal,
					Value: Str("0"),
				},
			},
		},
	})
}

func runTests(t *testing.T, testcases []TestCase) {

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			tokensStream := lexer.Lex(tc.code)
			tokens := []lexer.Tokval{}
			
			for t := range tokensStream {
				tokens = append(tokens, t)
			}
			
			assertEqualTokens(t, tc.want, tokens)
		})
	}
}

func assertEqualTokens(t *testing.T, want []lexer.Tokval, got []lexer.Tokval) {
	if len(want) != len(got) {
		t.Errorf("wanted [%d] tokens, got [%d] tokens", len(want), len(got))
		t.Fatalf("want[%v] != got[%v]", want, got)
	}
	
	for i, w := range want {
		g := got[i]
		if !w.Equal(g) {
			t.Errorf("wanted token[%s] != got token[%s]", w, g)
		}
	} 
}