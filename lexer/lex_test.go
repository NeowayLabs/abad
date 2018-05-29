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
			want: tokens(decimalToken("0")),
		},
		{
			name: "BigDecimal",
			code: Str("1236547987794465977"),
			want: tokens(decimalToken("1236547987794465977")),
		},
		{
			name: "RealDecimalStartingWithPoint",
			code: Str(".1"),
			want: tokens(decimalToken(".1")),
		},
		{
			name: "RealDecimalEndingWithPoint",
			code: Str("1."),
			want: tokens(decimalToken("1.")),
		},
		{
			name: "LargeRealDecimalStartingWithPoint",
			code: Str(".123456789"),
			want: tokens(decimalToken(".123456789")),
		},
		{
			name: "SmallRealDecimal",
			code: Str("1.6"),
			want: tokens(decimalToken("1.6")),
		},
		{
			name: "BigRealDecimal",
			code: Str("11223243554.63445465789"),
			want: tokens(decimalToken("11223243554.63445465789")),
		},
		{
			name: "SmallRealDecimalWithSmallExponent",
			code: Str("1.0e1"),
			want: tokens(decimalToken("1.0e1")),
		},
		{
			name: "SmallDecimalWithSmallExponent",
			code: Str("1e1"),
			want: tokens(decimalToken("1e1")),
		},
		{
			name: "SmallDecimalWithSmallExponentUpperExponent",
			code: Str("1E1"),
			want: tokens(decimalToken("1E1")),
		},
		{
			name: "BigDecimalWithBigExponent",
			code: Str("666666666666e668"),
			want: tokens(decimalToken("666666666666e668")),
		},
		{
			name: "BigDecimalWithBigExponentUpperExponent",
			code: Str("666666666666E668"),
			want: tokens(decimalToken("666666666666E668")),
		},
		{
			name: "BigRealDecimalWithBigExponent",
			code: Str("666666666666.0e66"),
			want: tokens(decimalToken("666666666666.0e66")),
		},
		{
			name: "RealDecimalWithSmallNegativeExponent",
			code: Str("1.0e-1"),
			want: tokens(decimalToken("1.0e-1")),
		},
		{
			name: "RealDecimalWithBigNegativeExponent",
			code: Str("1.0e-50"),
			want: tokens(decimalToken("1.0e-50")),
		},		
		{
			name: "SmallRealDecimalWithSmallUpperExponent",
			code: Str("1.0E1"),
			want: tokens(decimalToken("1.0E1")),
		},
		{
			name: "BigRealDecimalWithBigUpperExponent",
			code: Str("666666666666.0E66"),
			want: tokens(decimalToken("666666666666.0E66")),
		},
		{
			name: "RealDecimalWithSmallNegativeUpperExponent",
			code: Str("1.0E-1"),
			want: tokens(decimalToken("1.0E-1")),
		},
		{
			name: "RealDecimalWithBigNegativeUpperExponent",
			code: Str("1.0E-50"),
			want: tokens(decimalToken("1.0E-50")),
		},
		{
			name: "StartWithDotUpperExponent",
			code: Str(".0E-50"),
			want: tokens(decimalToken(".0E-50")),
		},
		{
			name: "StartWithDotExponent",
			code: Str(".0e5"),
			want: tokens(decimalToken(".0e5")),
		},
		{
			name: "ZeroHexadecimal",
			code: Str("0x0"),
			want: tokens(hexToken("0x0")),
		},
		{
			name: "BigHexadecimal",
			code: Str("0x123456789abcdef"),
			want: tokens(hexToken("0x123456789abcdef")),
		},
		{
			name: "BigHexadecimalUppercase",
			code: Str("0x123456789ABCDEF"),
			want: tokens(hexToken("0x123456789ABCDEF")),
		},
		{
			name: "LettersOnlyHexadecimal",
			code: Str("0xabcdef"),
			want: tokens(hexToken("0xabcdef")),
		},
		{
			name: "LettersOnlyHexadecimalUppercase",
			code: Str("0xABCDEF"),
			want: tokens(hexToken("0xABCDEF")),
		},
		{
			name: "ZeroHexadecimalUpperX",
			code: Str("0X0"),
			want: tokens(hexToken("0X0")),
		},
		{
			name: "BigHexadecimalUpperX",
			code: Str("0X123456789abcdef"),
			want: tokens(hexToken("0X123456789abcdef")),
		},
		{
			name: "BigHexadecimalUppercaseUpperX",
			code: Str("0X123456789ABCDEF"),
			want: tokens(hexToken("0X123456789ABCDEF")),
		},
		{
			name: "LettersOnlyHexadecimalUpperX",
			code: Str("0Xabcdef"),
			want: tokens(hexToken("0Xabcdef")),
		},
		{
			name: "LettersOnlyHexadecimalUppercaseUpperX",
			code: Str("0XABCDEF"),
			want: tokens(hexToken("0XABCDEF")),
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
		{
			name: "NumbersInterwined",
			code: Str("a1b2c"),
			want: tokens(identToken("a1b2c")),
		},
		{
			name: "AccessingMember",
			code: Str("console.log"),
			want: tokens(
				identToken("console"),
				dotToken(),
				identToken("log"),
			),
		},
		{
			name: "AccessingMemberOfMember",
			code: Str("console.log.toString"),
			want: tokens(
				identToken("console"),
				dotToken(),
				identToken("log"),
				dotToken(),
				identToken("toString"),
			),
		},
	})
}

func TestFuncall(t *testing.T) {
	runTests(t, []TestCase{
		{
			name: "OneLetterFunction",
			code: Str("a()"),
			want: tokens(
				identToken("a"),
				leftParenToken(),
				rightParenToken(),
			),
		},
		{
			name: "BigFunctionName",
			code: Str("veryBigFunctionNameThatWouldAnnoyNatel()"),
			want: tokens(
				identToken("veryBigFunctionNameThatWouldAnnoyNatel"),
				leftParenToken(),
				rightParenToken(),
			),
		},
		{
			name: "MemberFunction",
			code: Str("console.log()"),
			want: tokens(
				identToken("console"),
				dotToken(),
				identToken("log"),
				leftParenToken(),
				rightParenToken(),
			),
		},
		{
			name: "WithThreeDigitsDecimalParameter",
			code: Str("test(666)"),
			want: tokens(
				identToken("test"),
				leftParenToken(),
				decimalToken("666"),
				rightParenToken(),
			),
		},
		{
			name: "WithTwoDigitsDecimalParameter",
			code: Str("test(66)"),
			want: tokens(
				identToken("test"),
				leftParenToken(),
				decimalToken("66"),
				rightParenToken(),
			),
		},
		{
			name: "WithOneDigitDecimalParameter",
			code: Str("test(6)"),
			want: tokens(
				identToken("test"),
				leftParenToken(),
				decimalToken("6"),
				rightParenToken(),
			),
		},
		{
			name: "DecimalWithExponentParameter",
			code: Str("test(1e6)"),
			want: tokens(
				identToken("test"),
				leftParenToken(),
				decimalToken("1e6"),
				rightParenToken(),
			),
		},
		{
			name: "DecimalWithUpperExponentParameter",
			code: Str("test(1E6)"),
			want: tokens(
				identToken("test"),
				leftParenToken(),
				decimalToken("1E6"),
				rightParenToken(),
			),
		},
		{
			name: "WithSmallestRealDecimalParameter",
			code: Str("test(.1)"),
			want: tokens(
				identToken("test"),
				leftParenToken(),
				decimalToken(".1"),
				rightParenToken(),
			),
		},
		{
			name: "RealDecimalWithExponentParameter",
			code: Str("test(1.1e6)"),
			want: tokens(
				identToken("test"),
				leftParenToken(),
				decimalToken("1.1e6"),
				rightParenToken(),
			),
		},
		{
			name: "RealDecimalWithUpperExponentParameter",
			code: Str("test(1.1E6)"),
			want: tokens(
				identToken("test"),
				leftParenToken(),
				decimalToken("1.1E6"),
				rightParenToken(),
			),
		},
		{
			name: "WithRealDecimalParameter",
			code: Str("test(6.6)"),
			want: tokens(
				identToken("test"),
				leftParenToken(),
				decimalToken("6.6"),
				rightParenToken(),
			),
		},
		{
			name: "WithOneDigitHexadecimalParameter",
			code: Str("test(0x6)"),
			want: tokens(
				identToken("test"),
				leftParenToken(),
				hexToken("0x6"),
				rightParenToken(),
			),
		},
		{
			name: "WithOneDigitUpperHexadecimalParameter",
			code: Str("test(0X6)"),
			want: tokens(
				identToken("test"),
				leftParenToken(),
				hexToken("0X6"),
				rightParenToken(),
			),
		},
		{
			name: "WithTwoDigitHexadecimalParameter",
			code: Str("test(0x66)"),
			want: tokens(
				identToken("test"),
				leftParenToken(),
				hexToken("0x66"),
				rightParenToken(),
			),
		},
		{
			name: "WithTwoDigitUpperHexadecimalParameter",
			code: Str("test(0X66)"),
			want: tokens(
				identToken("test"),
				leftParenToken(),
				hexToken("0X66"),
				rightParenToken(),
			),
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

func TestIllegalMemberAccess(t *testing.T) {

	runTests(t, []TestCase{
		{
			name: "CantAccessMemberThatStartsWithNumber",
			code: Str("test.123"),
			want: []lexer.Tokval{
				identToken("test"),
				dotToken(),
				illegalToken("123"),
			},
		},
		{
			name: "CantAccessMemberThatStartsWithDot",
			code: Str("test.."),
			want: []lexer.Tokval{
				identToken("test"),
				dotToken(),
				illegalToken("."),
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
				illegalToken("123E123E123"),
			},
		},
		{
			name: "DecimalDuplicatedExponentPart",
			code: Str("123e123e123"),
			want: []lexer.Tokval{
				illegalToken("123e123e123"),
			},
		},
		{
			name: "RealDecimalDuplicatedUpperExponentPart",
			code: Str("123.1E123E123"),
			want: []lexer.Tokval{
				illegalToken("123.1E123E123"),
			},
		},
		{
			name: "RealDecimalDuplicatedExponentPart",
			code: Str("123.6e123e123"),
			want: []lexer.Tokval{
				illegalToken("123.6e123e123"),
			},
		},
		{
			name: "OnlyStartAsDecimal",
			code: Str("0LALALA"),
			want: []lexer.Tokval{
				illegalToken("0LALALA"),
			},
		},
		{
			name: "EndIsNotDecimal",
			code: Str("0123344546I4K"),
			want: []lexer.Tokval{
				illegalToken("0123344546I4K"),
			},
		},
		{
			name: "EmptyHexadecimal",
			code: Str("0x"),
			want: []lexer.Tokval{
				illegalToken("0x"),
			},
		},
		{
			name: "OnlyStartAsReal",
			code: Str("0.b"),
			want: []lexer.Tokval{
				illegalToken("0.b"),
			},
		},
		{
			name: "RealWithTwoDotsStartingWithDot",
			code: Str(".1.2"),
			want: []lexer.Tokval{
				illegalToken(".1.2"),
			},
		},
		{
			name: "RealWithTwoDots",
			code: Str("0.1.2"),
			want: []lexer.Tokval{
				illegalToken("0.1.2"),
			},
		},
		{
			name: "BifRealWithTwoDots",
			code: Str("1234.666.2342"),
			want: []lexer.Tokval{
				illegalToken("1234.666.2342"),
			},
		},
		{
			name: "EmptyHexadecimalUpperX",
			code: Str("0X"),
			want: []lexer.Tokval{
				illegalToken("0X"),
			},
		},
		{
			name: "LikeHexadecimal",
			code: Str("0b1234"),
			want: []lexer.Tokval{
				illegalToken("0b1234"),
			},
		},
		{
			name: "OnlyStartAsHexadecimal",
			code: Str("0xI4K"),
			want: []lexer.Tokval{
				illegalToken("0xI4K"),
			},
		},
		{
			name: "EndIsNotHexadecimal",
			code: Str("0x123456G"),
			want: []lexer.Tokval{
				illegalToken("0x123456G"),
			},
		},
		{
			name: "CorruptedHexadecimal",
			code: corruptedHex,
			want: []lexer.Tokval{
				illegalToken(corruptedHex.String()),
			},
		},
		{
			name: "CorruptedDecimal",
			code: corruptedDecimal,
			want: []lexer.Tokval{
				illegalToken(corruptedDecimal.String()),
			},
		},
		{
			name: "CorruptedNumber",
			code: corruptedNumber,
			want: []lexer.Tokval{
				illegalToken(corruptedNumber.String()),
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
			want: []lexer.Tokval{ illegalToken(messStr(Str(""), 0).String()) },
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

func illegalToken(val string) lexer.Tokval {
	return lexer.Tokval{
		Type: token.Illegal,
		Value: Str(val),
	}
}

func assertWantedTokens(t *testing.T, tc TestCase, got []lexer.Tokval) {
	t.Helper()
	
	if len(tc.want) != len(got) {
		t.Errorf("wanted [%d] tokens, got [%d] tokens", len(tc.want), len(got))
		t.Fatalf("\nwant=%v\ngot= %v\nare not equal.", tc.want, got)
	}
	
	for i, w := range tc.want {
		g := got[i]
		if !w.Equal(g) {
			t.Errorf("wanted token[%d][%+v] != got token[%d][%+v]", i, w, i, g)
			t.Errorf("wanted tokens[%+v] != got tokens[%+v]", tc.want, got)
		}
		
		if tc.checkPosition {
			if !w.EqualPos(g) {
				t.Errorf("want=%+v\ngot=%+v\nare equal but dont have the same position",w, g)
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

func leftParenToken() lexer.Tokval {
	return lexer.Tokval{
		Type: token.LParen,
		Value: Str("("),
	}
}

func rightParenToken() lexer.Tokval {
	return lexer.Tokval{
		Type: token.RParen,
		Value: Str(")"),
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

func dotToken() lexer.Tokval {
	return lexer.Tokval{
		Type: token.Dot,
		Value: Str("."),
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