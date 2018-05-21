package lexer_test

import (
	"testing"
	
	"github.com/NeowayLabs/abad/token"
	"github.com/NeowayLabs/abad/lexer"
	"github.com/NeowayLabs/abad/internal/utf16"
)

type Token struct {
	Type token.Type
	Value string
}

type TestCase struct {
	name string
	code string
	want []Token
}

func TestNumericLiterals(t *testing.T) {

	runTests(t, []TestCase{
		{
			name: "Zero",
			code: "0",
			want: []Token{
				{
					Type: token.Decimal,
					Value: "0",
				},
			},
		},
	})
}

func runTests(t *testing.T, testcases []TestCase) {

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			tokensStream := lexer.Lex(utf16.S(tc.code))
			tokens := []lexer.Tokval{}
			
			for t := range tokensStream {
				tokens = append(tokens, t)
			}
			
			assertEqualTokens(t, tokvals(tc.want), tokens)
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

func tokvals(tokens []Token) []lexer.Tokval {
	tokvals := make([]lexer.Tokval, len(tokens))
	
	for i, t := range tokens {
		tokvals[i] = lexer.Tokval{
			Type: t.Type,
			Value: utf16.S(t.Value),
		}
	}
	
	return tokvals
}