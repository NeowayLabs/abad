package lexer_test

import (
	"fmt"
	"testing"
	"unicode"
	
	"github.com/NeowayLabs/abad/token"
	"github.com/NeowayLabs/abad/lexer"
	"github.com/NeowayLabs/abad/internal/utf16"
)


type TestCase struct {
	name string
	code utf16.Str
	want []lexer.Tokval
	checkPosition bool
}

var Str func(string) utf16.Str = utf16.S
var EOF lexer.Tokval = lexer.EOF

func TestNumericLiterals(t *testing.T) {

	// SPEC: https://es5.github.io/#x7.8.3
	
	cases := []TestCase{
		{
			name: "SingleZero",
			code: Str("0"),
			want: []lexer.Tokval{
				decimalToken("0"),
				EOF,
			},
		},
		{
			name: "BigDecimal",
			code: Str("1236547987794465977"),
			want: []lexer.Tokval{
				decimalToken("1236547987794465977"),
				EOF,
			},
		},
		{
			name: "RealDecimalStartingWithPoint",
			code: Str(".1"),
			want: []lexer.Tokval{
				decimalToken(".1"),
				EOF,
			},
		},
		{
			name: "RealDecimalEndingWithPoint",
			code: Str("1."),
			want: []lexer.Tokval{
				decimalToken("1."),
				EOF,
			},
		},
		{
			name: "LargeRealDecimalStartingWithPoint",
			code: Str(".123456789"),
			want: []lexer.Tokval{
				decimalToken(".123456789"),
				EOF,
			},
		},
		{
			name: "SmallRealDecimal",
			code: Str("1.6"),
			want: []lexer.Tokval{
				decimalToken("1.6"),
				EOF,
			},
		},
		{
			name: "BigRealDecimal",
			code: Str("11223243554.63445465789"),
			want: []lexer.Tokval{
				decimalToken("11223243554.63445465789"),
				EOF,
			},
		},
		{
			name: "SmallRealDecimalWithSmallExponent",
			code: Str("1.0e1"),
			want: []lexer.Tokval{
				decimalToken("1.0e1"),
				EOF,
			},
		},
		{
			name: "SmallDecimalWithSmallExponent",
			code: Str("1e1"),
			want: []lexer.Tokval{
				decimalToken("1e1"),
				EOF,
			},
		},
		{
			name: "SmallDecimalWithSmallExponentUpperExponent",
			code: Str("1E1"),
			want: []lexer.Tokval{
				decimalToken("1E1"),
				EOF,
			},
		},
		{
			name: "BigDecimalWithBigExponent",
			code: Str("666666666666e668"),
			want: []lexer.Tokval{
				decimalToken("666666666666e668"),
				EOF,
			},
		},
		{
			name: "BigDecimalWithBigExponentUpperExponent",
			code: Str("666666666666E668"),
			want: []lexer.Tokval{
				decimalToken("666666666666E668"),
				EOF,
			},
		},
		{
			name: "BigRealDecimalWithBigExponent",
			code: Str("666666666666.0e66"),
			want: []lexer.Tokval{
				decimalToken("666666666666.0e66"),
				EOF,
			},
		},
		{
			name: "RealDecimalWithSmallNegativeExponent",
			code: Str("1.0e-1"),
			want: []lexer.Tokval{
				decimalToken("1.0e-1"),
				EOF,
			},
		},
		{
			name: "RealDecimalWithBigNegativeExponent",
			code: Str("1.0e-50"),
			want: []lexer.Tokval{
				decimalToken("1.0e-50"),
				EOF,
			},
		},		
		{
			name: "SmallRealDecimalWithSmallUpperExponent",
			code: Str("1.0E1"),
			want: []lexer.Tokval{
				decimalToken("1.0E1"),
				EOF,
			},
		},
		{
			name: "BigRealDecimalWithBigUpperExponent",
			code: Str("666666666666.0E66"),
			want: []lexer.Tokval{
				decimalToken("666666666666.0E66"),
				EOF,
			},
		},
		{
			name: "RealDecimalWithSmallNegativeUpperExponent",
			code: Str("1.0E-1"),
			want: []lexer.Tokval{
				decimalToken("1.0E-1"),
				EOF,
			},
		},
		{
			name: "RealDecimalWithBigNegativeUpperExponent",
			code: Str("1.0E-50"),
			want: []lexer.Tokval{
				decimalToken("1.0E-50"),
				EOF,
			},
		},
		{
			name: "StartWithDotUpperExponent",
			code: Str(".0E-50"),
			want: []lexer.Tokval{
				decimalToken(".0E-50"),
				EOF,
			},
		},
		{
			name: "StartWithDotExponent",
			code: Str(".0e5"),
			want: []lexer.Tokval{
				decimalToken(".0e5"),
				EOF,
			},
		},
		{
			name: "ZeroHexadecimal",
			code: Str("0x0"),
			want: []lexer.Tokval{
				hexToken("0x0"),
				EOF,
			},
		},
		{
			name: "BigHexadecimal",
			code: Str("0x123456789abcdef"),
			want: []lexer.Tokval{
				hexToken("0x123456789abcdef"),
				EOF,
			},
		},
		{
			name: "BigHexadecimalUppercase",
			code: Str("0x123456789ABCDEF"),
			want: []lexer.Tokval{
				hexToken("0x123456789ABCDEF"),
				EOF,
			},
		},
		{
			name: "LettersOnlyHexadecimal",
			code: Str("0xabcdef"),
			want: []lexer.Tokval{
				hexToken("0xabcdef"),
				EOF,
			},
		},
		{
			name: "LettersOnlyHexadecimalUppercase",
			code: Str("0xABCDEF"),
			want: []lexer.Tokval{
				hexToken("0xABCDEF"),
				EOF,
			},
		},
		{
			name: "ZeroHexadecimalUpperX",
			code: Str("0X0"),
			want: []lexer.Tokval{
				hexToken("0X0"),
				EOF,
			},
		},
		{
			name: "BigHexadecimalUpperX",
			code: Str("0X123456789abcdef"),
			want: []lexer.Tokval{
				hexToken("0X123456789abcdef"),
				EOF,
			},
		},
		{
			name: "BigHexadecimalUppercaseUpperX",
			code: Str("0X123456789ABCDEF"),
			want: []lexer.Tokval{
				hexToken("0X123456789ABCDEF"),
				EOF,
			},
		},
		{
			name: "LettersOnlyHexadecimalUpperX",
			code: Str("0Xabcdef"),
			want: []lexer.Tokval{
				hexToken("0Xabcdef"),
				EOF,
			},
		},
		{
			name: "LettersOnlyHexadecimalUppercaseUpperX",
			code: Str("0XABCDEF"),
			want: []lexer.Tokval{
				hexToken("0XABCDEF"),
				EOF,
			},
		},
	}
	
	plusSignedCases := prependOnTestCases(TestCase{
		name: "PlusSign",
		code: Str("+"),
		want: []lexer.Tokval{ plusToken() },
	}, cases)
	
	minusSignedCases := prependOnTestCases(TestCase{
		name: "MinusSign",
		code: Str("-"),
		want: []lexer.Tokval{ minusToken() },
	}, cases)
	
	plusMinusPlusMinusSignedCases := prependOnTestCases(TestCase{
		name: "PlusMinusPlusMinusSign",
		code: Str("+-+-"),
		want: []lexer.Tokval{
			plusToken(),
			minusToken(),
			plusToken(),
			minusToken(),
		},
	}, cases)
	
	minusPlusMinusPlusSignedCases := prependOnTestCases(TestCase{
		name: "MinusPlusMinusPlusSign",
		code: Str("-+-+"),
		want: []lexer.Tokval{
			minusToken(),
			plusToken(),
			minusToken(),
			plusToken(),
		},
	}, cases)
	
	runTests(t, cases)
	runTests(t, plusSignedCases)
	runTests(t, minusSignedCases)
	runTests(t, plusMinusPlusMinusSignedCases)
	runTests(t, minusPlusMinusPlusSignedCases)
}

func TestIdentifiers(t *testing.T) {
	runTests(t, []TestCase{
		{
			name: "Underscore",
			code: Str("_"),
			want: tokens(identToken("_")),
		},
		{
			name: "SingleLetter",
			code: Str("a"),
			want: tokens(identToken("a")),
		},
		{
			name: "Self",
			code: Str("self"),
			want: tokens(identToken("self")),
		},
		{
			name: "Console",
			code: Str("console"),
			want: tokens(identToken("console")),
		},
		{
			name: "LotsUnderscores",
			code: Str("___hyped___"),
			want: tokens(identToken("___hyped___")),
		},
		{
			name: "DollarsInterwined",
			code: Str("a$b$c"),
			want: tokens(identToken("a$b$c")),
		},
	})
}

func TestPosition(t *testing.T) {
	runTests(t, []TestCase{
		{
			name: "MinusDecimal",
			code: Str("-1"),
			checkPosition: true,
			want: []lexer.Tokval{
				minusTokenPos(1,1),
				decimalTokenPos("1", 1, 2),
				EOF,
			},
		},
		{
			name: "PlusDecimal",
			code: Str("+1"),
			checkPosition: true,
			want: []lexer.Tokval{
				plusTokenPos(1,1),
				decimalTokenPos("1", 1, 2),
				EOF,
			},
		},
		{
			name: "PlusMinusDecimal",
			code: Str("+-666"),
			checkPosition: true,
			want: []lexer.Tokval{
				plusTokenPos(1,1),
				minusTokenPos(1,2),
				decimalTokenPos("666", 1, 3),
				EOF,
			},
		},
	})
}

func TestIllegalNumericLiterals(t *testing.T) {
	
	corruptedHex := messStr(Str("0x01234"), 4)
	corruptedDecimal := messStr(Str("1234"), 3)
	corruptedNumber := messStr(Str("0"), 1)
	
	runTests(t, []TestCase{
		{
			name: "DecimalDuplicatedUpperExponentPart",
			code: Str("123E123E123"),
			want: []lexer.Tokval{
				illegalToken(Str("123E123E123")),
			},
		},
		{
			name: "DecimalDuplicatedExponentPart",
			code: Str("123e123e123"),
			want: []lexer.Tokval{
				illegalToken(Str("123e123e123")),
			},
		},
		{
			name: "RealDecimalDuplicatedUpperExponentPart",
			code: Str("123.1E123E123"),
			want: []lexer.Tokval{
				illegalToken(Str("123.1E123E123")),
			},
		},
		{
			name: "RealDecimalDuplicatedExponentPart",
			code: Str("123.6e123e123"),
			want: []lexer.Tokval{
				illegalToken(Str("123.6e123e123")),
			},
		},
		{
			name: "OnlyStartAsDecimal",
			code: Str("0LALALA"),
			want: []lexer.Tokval{
				illegalToken(Str("0LALALA")),
			},
		},
		{
			name: "EndIsNotDecimal",
			code: Str("0123344546I4K"),
			want: []lexer.Tokval{
				illegalToken(Str("0123344546I4K")),
			},
		},
		{
			name: "EmptyHexadecimal",
			code: Str("0x"),
			want: []lexer.Tokval{
				illegalToken(Str("0x")),
			},
		},
		{
			name: "OnlyStartAsReal",
			code: Str("0.b"),
			want: []lexer.Tokval{
				illegalToken(Str("0.b")),
			},
		},
		{
			name: "RealWithTwoDotsStartingWithDot",
			code: Str(".1.2"),
			want: []lexer.Tokval{
				illegalToken(Str(".1.2")),
			},
		},
		{
			name: "RealWithTwoDots",
			code: Str("0.1.2"),
			want: []lexer.Tokval{
				illegalToken(Str("0.1.2")),
			},
		},
		{
			name: "BifRealWithTwoDots",
			code: Str("1234.666.2342"),
			want: []lexer.Tokval{
				illegalToken(Str("1234.666.2342")),
			},
		},
		{
			name: "EmptyHexadecimalUpperX",
			code: Str("0X"),
			want: []lexer.Tokval{
				illegalToken(Str("0X")),
			},
		},
		{
			name: "LikeHexadecimal",
			code: Str("0b1234"),
			want: []lexer.Tokval{
				illegalToken(Str("0b1234")),
			},
		},
		{
			name: "OnlyStartAsHexadecimal",
			code: Str("0xI4K"),
			want: []lexer.Tokval{
				illegalToken(Str("0xI4K")),
			},
		},
		{
			name: "EndIsNotHexadecimal",
			code: Str("0x123456G"),
			want: []lexer.Tokval{
				illegalToken(Str("0x123456G")),
			},
		},
		{
			name: "CorruptedHexadecimal",
			code: corruptedHex,
			want: []lexer.Tokval{
				illegalToken(corruptedHex),
			},
		},
		{
			name: "CorruptedDecimal",
			code: corruptedDecimal,
			want: []lexer.Tokval{
				illegalToken(corruptedDecimal),
			},
		},
		{
			name: "CorruptedNumber",
			code: corruptedNumber,
			want: []lexer.Tokval{
				illegalToken(corruptedNumber),
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
	})
}

func TestCorruptedInput(t *testing.T) {
	runTests(t, []TestCase{
		{
			name: "AtStart",
			code: messStr(Str(""), 0),
			want: []lexer.Tokval{ illegalToken(messStr(Str(""), 0)) },
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
			
			assertWantedTokens(t, tc, tokens)
		})
	}
}

func illegalToken(val utf16.Str) lexer.Tokval {
	return lexer.Tokval{
		Type: token.Illegal,
		Value: val,
	}
}

func assertWantedTokens(t *testing.T, tc TestCase, got []lexer.Tokval) {
	t.Helper()
	
	if len(tc.want) != len(got) {
		t.Errorf("wanted [%d] tokens, got [%d] tokens", len(tc.want), len(got))
		t.Fatalf("want[%+v] != got[%+v]", tc.want, got)
	}
	
	for i, w := range tc.want {
		g := got[i]
		if !w.Equal(g) {
			t.Errorf("wanted token[%d][%+v] != got token[%d][%+v]", i, w, i, g)
			t.Errorf("wanted tokens[%+v] != got tokens[%+v]", tc.want, got)
		}
		
		if tc.checkPosition {
			if !w.EqualPos(g) {
				t.Errorf("want[%+v] got[%+v] are equal but dont have the same position",w, g)
			}
		}
	} 
}

func messStr(s utf16.Str, pos uint) utf16.Str {
	// WHY: The go's utf16 package uses the replacement char everytime a some
	// encoding/decoding error happens, so we inject one on the uint16 array to simulate
	// encoding/decoding errors.
	// Not safe but the idea is to fuck up the string	

	r := append(s[0:pos], uint16(unicode.ReplacementChar))
	r = append(r, s[pos:]...)
	return r
}

// prependOnTestCases will prepend the given tcase on each TestCase
// on tcases generating a new array of TestCases.
//
// The array of TestCases is generated by prepending code and the
// wanted tokens for each TestCases. EOF should not be provided on the
// given tcase since it will be prepended on each test case inside given tcases.
func prependOnTestCases(tcase TestCase, tcases []TestCase) []TestCase{
	newcases := make([]TestCase, len(tcases))
	
	for i, t := range tcases {
		name := fmt.Sprintf("%s/%s", tcase.name, t.name)
		code := tcase.code.Append(t.code)
		want := append(tcase.want, t.want...)
		
		newcases[i] = TestCase{
			name: name,
			code: code,
			want: want,
		}
	}
	
	return newcases
}

func minusToken() lexer.Tokval {
	return lexer.Tokval{
		Type: token.Minus,
		Value: Str("-"),
	}
}

func plusToken() lexer.Tokval {
	return lexer.Tokval{
		Type: token.Plus,
		Value: Str("+"),
	}
}

func minusTokenPos(line uint, column uint) lexer.Tokval {
	return lexer.Tokval{
		Type: token.Minus,
		Value: Str("-"),
		Line: line,
		Column: column,
	}
}

func plusTokenPos(line uint, column uint) lexer.Tokval {
	return lexer.Tokval{
		Type: token.Plus,
		Value: Str("+"),
		Line: line,
		Column: column,
	}
}

func decimalTokenPos(dec string, line uint, column uint) lexer.Tokval {
	return lexer.Tokval{
		Type: token.Decimal,
		Value: Str(dec),
		Line: line,
		Column: column,
	}
}

func decimalToken(dec string) lexer.Tokval {
	return lexer.Tokval{
		Type: token.Decimal,
		Value: Str(dec),
	}
}

func hexToken(hex string) lexer.Tokval {
	return lexer.Tokval{
		Type: token.Hexadecimal,
		Value: Str(hex),
	}
}

func identToken(s string) lexer.Tokval {
	return lexer.Tokval{
		Type: token.Ident,
		Value: Str(s),
	}
}

func tokens(t ...lexer.Tokval) []lexer.Tokval {
	return append(t, EOF)
}