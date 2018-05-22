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

	// SPEC: https://es5.github.io/#x7.8.3
	
	runTests(t, []TestCase{
		{
			name: "SingleZero",
			code: Str("0"),
			want: []lexer.Tokval{
				{
					Type: token.Decimal,
					Value: Str("0"),
				},
			},
		},
		{
			name: "BigDecimal",
			code: Str("1236547987794465977"),
			want: []lexer.Tokval{
				{
					Type: token.Decimal,
					Value: Str("1236547987794465977"),
				},
			},
		},
		{
			name: "RealDecimalStartingWithPoint",
			code: Str(".1"),
			want: []lexer.Tokval{
				{
					Type: token.Decimal,
					Value: Str(".1"),
				},
			},
		},
		{
			name: "RealDecimalEndingWithPoint",
			code: Str("1."),
			want: []lexer.Tokval{
				{
					Type: token.Decimal,
					Value: Str("1."),
				},
			},
		},
		{
			name: "LargeRealDecimalStartingWithPoint",
			code: Str(".123456789"),
			want: []lexer.Tokval{
				{
					Type: token.Decimal,
					Value: Str(".123456789"),
				},
			},
		},
		{
			name: "SmallRealDecimal",
			code: Str("1.6"),
			want: []lexer.Tokval{
				{
					Type: token.Decimal,
					Value: Str("1.6"),
				},
			},
		},
		{
			name: "BigRealDecimal",
			code: Str("11223243554.63445465789"),
			want: []lexer.Tokval{
				{
					Type: token.Decimal,
					Value: Str("11223243554.63445465789"),
				},
			},
		},
		{
			name: "SmallRealDecimalWithSmallExponent",
			code: Str("1.0e1"),
			want: []lexer.Tokval{
				{
					Type: token.Decimal,
					Value: Str("1.0e1"),
				},
			},
		},
		{
			name: "BigRealDecimalWithBigExponent",
			code: Str("666666666666.0e66"),
			want: []lexer.Tokval{
				{
					Type: token.Decimal,
					Value: Str("666666666666.0e66"),
				},
			},
		},
		{
			name: "RealDecimalWithSmallNegativeExponent",
			code: Str("1.0e-1"),
			want: []lexer.Tokval{
				{
					Type: token.Decimal,
					Value: Str("1.0e-1"),
				},
			},
		},
		{
			name: "RealDecimalWithBigNegativeExponent",
			code: Str("1.0e-50"),
			want: []lexer.Tokval{
				{
					Type: token.Decimal,
					Value: Str("1.0e-50"),
				},
			},
		},		
		{
			name: "SmallRealDecimalWithSmallUpperExponent",
			code: Str("1.0E1"),
			want: []lexer.Tokval{
				{
					Type: token.Decimal,
					Value: Str("1.0E1"),
				},
			},
		},
		{
			name: "BigRealDecimalWithBigUpperExponent",
			code: Str("666666666666.0E66"),
			want: []lexer.Tokval{
				{
					Type: token.Decimal,
					Value: Str("666666666666.0E66"),
				},
			},
		},
		{
			name: "RealDecimalWithSmallNegativeUpperExponent",
			code: Str("1.0E-1"),
			want: []lexer.Tokval{
				{
					Type: token.Decimal,
					Value: Str("1.0E-1"),
				},
			},
		},
		{
			name: "RealDecimalWithBigNegativeUpperExponent",
			code: Str("1.0E-50"),
			want: []lexer.Tokval{
				{
					Type: token.Decimal,
					Value: Str("1.0E-50"),
				},
			},
		},
	})
}

func TestIllegalNumericLiterals(t *testing.T) {
	// TODO
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