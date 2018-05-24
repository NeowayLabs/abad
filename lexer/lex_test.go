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
var EOF lexer.Tokval = lexer.Tokval{ Type: token.EOF }

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
				EOF,
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
				EOF,
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
				EOF,
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
				EOF,
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
				EOF,
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
				EOF,
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
				EOF,
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
				EOF,
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
				EOF,
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
				EOF,
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
				EOF,
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
				EOF,
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
				EOF,
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
				EOF,
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
				EOF,
			},
		},
		{
			name: "ZeroHexadecimal",
			code: Str("0x0"),
			want: []lexer.Tokval{
				{
					Type: token.Hexadecimal,
					Value: Str("0x0"),
				},
				EOF,
			},
		},
		{
			name: "BigHexadecimal",
			code: Str("0x123456789abcdef"),
			want: []lexer.Tokval{
				{
					Type: token.Hexadecimal,
					Value: Str("0x123456789abcdef"),
				},
				EOF,
			},
		},
		{
			name: "BigHexadecimalUppercase",
			code: Str("0x123456789ABCDEF"),
			want: []lexer.Tokval{
				{
					Type: token.Hexadecimal,
					Value: Str("0x123456789ABCDEF"),
				},
				EOF,
			},
		},
		{
			name: "LettersOnlyHexadecimal",
			code: Str("0xabcdef"),
			want: []lexer.Tokval{
				{
					Type: token.Hexadecimal,
					Value: Str("0xabcdef"),
				},
				EOF,
			},
		},
		{
			name: "LettersOnlyHexadecimalUppercase",
			code: Str("0xABCDEF"),
			want: []lexer.Tokval{
				{
					Type: token.Hexadecimal,
					Value: Str("0xABCDEF"),
				},
				EOF,
			},
		},
		{
			name: "ZeroHexadecimalUpperX",
			code: Str("0X0"),
			want: []lexer.Tokval{
				{
					Type: token.Hexadecimal,
					Value: Str("0X0"),
				},
				EOF,
			},
		},
		{
			name: "BigHexadecimalUpperX",
			code: Str("0X123456789abcdef"),
			want: []lexer.Tokval{
				{
					Type: token.Hexadecimal,
					Value: Str("0X123456789abcdef"),
				},
				EOF,
			},
		},
		{
			name: "BigHexadecimalUppercaseUpperX",
			code: Str("0X123456789ABCDEF"),
			want: []lexer.Tokval{
				{
					Type: token.Hexadecimal,
					Value: Str("0X123456789ABCDEF"),
				},
				EOF,
			},
		},
		{
			name: "LettersOnlyHexadecimalUpperX",
			code: Str("0Xabcdef"),
			want: []lexer.Tokval{
				{
					Type: token.Hexadecimal,
					Value: Str("0Xabcdef"),
				},
				EOF,
			},
		},
		{
			name: "LettersOnlyHexadecimalUppercaseUpperX",
			code: Str("0XABCDEF"),
			want: []lexer.Tokval{
				{
					Type: token.Hexadecimal,
					Value: Str("0XABCDEF"),
				},
				EOF,
			},
		},
	})
}

func TestIllegalNumericLiterals(t *testing.T) {
	runTests(t, []TestCase{
		{
			name: "EmptyHexadecimal",
			code: Str("0x"),
			want: []lexer.Tokval{
				{
					Type: token.Illegal,
					Value: Str("0x"),
				},
			},
		},
		{
			name: "EmptyHexadecimalUpperX",
			code: Str("0X"),
			want: []lexer.Tokval{
				{
					Type: token.Illegal,
					Value: Str("0X"),
				},
			},
		},
		{
			name: "LikeHexadecimal",
			code: Str("0b1234"),
			want: []lexer.Tokval{
				{
					Type: token.Illegal,
					Value: Str("0b1234"),
				},
			},
		},
	})
}

func TestNoOutputFor(t *testing.T) {
	runTests(t, []TestCase{
		{
			name: "EmptyString",
			code: Str(""),
			want: []lexer.Tokval{ EOF },
		},
		{
			name: "JustSpaces",
			code: Str("        "),
			want: []lexer.Tokval{ EOF },
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
	t.Helper()
	
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