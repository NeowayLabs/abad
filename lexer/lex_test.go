package lexer_test

import (
	"fmt"
	"testing"
	"unicode"

	"github.com/NeowayLabs/abad/internal/utf16"
	"github.com/NeowayLabs/abad/lexer"
	"github.com/NeowayLabs/abad/token"
)

type TestCase struct {
	name          string
	code          utf16.Str
	want          []lexer.Tokval
	checkPosition bool
}

func (tc TestCase) String() string {
	return fmt.Sprintf("name[%s] code[%s] want[%v] checkPosition[%t]", tc.name, tc.code, tc.want, tc.checkPosition)
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
		want: []lexer.Tokval{plusToken()},
	}, cases)

	minusSignedCases := prependOnTestCases(TestCase{
		name: "MinusSign",
		code: Str("-"),
		want: []lexer.Tokval{minusToken()},
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

	cases = append(cases, plusSignedCases...)
	cases = append(cases, minusSignedCases...)
	cases = append(cases, plusMinusPlusMinusSignedCases...)
	cases = append(cases, minusPlusMinusPlusSignedCases...)

	runTests(t, cases)
	runTokenSepTests(t, cases)
}

func TestStrings(t *testing.T) {
	// TODO: multiline strings
	// - escaped double quotes

	cases := []TestCase{
		{
			name: "Empty",
			code: Str(`""`),
			want: tokens(stringToken("")),
		},
		{
			name: "SpacesOnly",
			code: Str(`"  "`),
			want: tokens(stringToken("  ")),
		},
		{
			name: "semicolon",
			code: Str(`";"`),
			want: tokens(stringToken(";")),
		},
		{
			name: "SingleChar",
			code: Str(`"k"`),
			want: tokens(stringToken("k")),
		},
		{
			name: "LotsOfCrap",
			code: Str(`"1234567890-+=abcdefg${[]})(()%_ /|/ yay %xi4klindaum"`),
			want: tokens(stringToken("1234567890-+=abcdefg${[]})(()%_ /|/ yay %xi4klindaum")),
		},
	}

	runTests(t, cases)
	runTokenSepTests(t, cases)
}

func TestKeywords(t *testing.T) {
	cases := []TestCase {
		{
			name: "Null",
			code: Str("null"),
			want: tokens(nullToken()),
		},
	}
	
	runTests(t, cases)
	runTokenSepTests(t, cases)
}

func TestSemiColon(t *testing.T) {
	// Almost all semicolon tests are made interwined on other tests
	runTests(t, []TestCase{
		{
			name: "SingleSemiColon",
			code: Str(";"),
			want: tokens(semiColonToken()),
		},
		{
			name: "MultipleSemiColon",
			code: Str(";;;"),
			want: tokens(semiColonToken(), semiColonToken(), semiColonToken()),
		},
	})
}

func TestLineTerminator(t *testing.T) {

	for _, lineTerminator := range lineTerminators() {
		t.Run(lineTerminator.name, func(t *testing.T) {
			lt := lineTerminator.val
			lttok := ltToken(lt)
			runTests(t, []TestCase{
				{
					name: fmt.Sprintf("Single%s", lt),
					code: Str(lt),
					want: tokens(lttok),
				},
				{
					name: "Strings",
					code: sfmt(`"first"%s"second"`, lt),
					want: tokens(stringToken("first"), lttok, stringToken("second")),
				},
				{
					name: "Decimals",
					code: sfmt("1%s2", lt),
					want: tokens(decimalToken("1"), lttok, decimalToken("2")),
				},
				{
					name: "ExponentDecimals",
					code: sfmt("1e1%s1e+1%s1e-1%s1", lt, lt, lt),
					want: tokens(
						decimalToken("1e1"),
						lttok,
						decimalToken("1e+1"),
						lttok,
						decimalToken("1e-1"),
						lttok,
						decimalToken("1"),
					),
				},
				{
					name: "RealDecimals",
					code: sfmt(".1%s245.123", lt),
					want: tokens(decimalToken(".1"), lttok, decimalToken("245.123")),
				},
				{
					name: "Hexadecimals",
					code: sfmt("0xFF%s0x11", lt),
					want: tokens(hexToken("0xFF"), lttok, hexToken("0x11")),
				},
				{
					name: "Identifiers",
					code: sfmt("hi%shello", lt),
					want: tokens(identToken("hi"), lttok, identToken("hello")),
				},
				{
					name: "TwoFuncalls",
					code: sfmt("func1(a)%sfunc2(1)%s", lt, lt),
					want: tokens(
						identToken("func1"),
						leftParenToken(),
						identToken("a"),
						rightParenToken(),
						lttok,
						identToken("func2"),
						leftParenToken(),
						decimalToken("1"),
						rightParenToken(),
						lttok,
					),
				},
				{
					name: "FuncallWithSemiColon",
					code: sfmt("a();%sb()", lt),
					want: tokens(
						identToken("a"),
						leftParenToken(),
						rightParenToken(),
						semiColonToken(),
						lttok,
						identToken("b"),
						leftParenToken(),
						rightParenToken(),
					),
				},
			})
		})
	}
}

func TestIdentifiers(t *testing.T) {

	identCases := []TestCase{
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
			name: "DollarsIntertwined",
			code: Str("a$b$c"),
			want: tokens(identToken("a$b$c")),
		},
		{
			name: "NumbersIntertwined",
			code: Str("a1b2c"),
			want: tokens(identToken("a1b2c")),
		},
	}

	accessModCases := []TestCase{
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
			name: "AccessingNoMember",
			code: Str("console."),
			want: tokens(
				identToken("console"),
				dotToken(),
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
	}

	runTests(t, identCases)
	runTests(t, accessModCases)

	runTokenSepTests(t, identCases)
}

func TestFuncall(t *testing.T) {
	// TODO: add anon funcall "(function (a) { console.log(a); })("hi")"
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
			name: "NestingWithSemiColonNewLine",
			code: Str("a(b(),c(d(),e(),5));\n"),
			want: tokens(
				identToken("a"),
				leftParenToken(),
				identToken("b"),
				leftParenToken(),
				rightParenToken(),
				commaToken(),
				identToken("c"),
				leftParenToken(),
				identToken("d"),
				leftParenToken(),
				rightParenToken(),
				commaToken(),
				identToken("e"),
				leftParenToken(),
				rightParenToken(),
				commaToken(),
				decimalToken("5"),
				rightParenToken(),
				rightParenToken(),
				semiColonToken(),
				ltToken("\n"),
			),
		},
		{
			name: "SeparatedBySemiColon",
			code: Str("a();b()"),
			want: tokens(
				identToken("a"),
				leftParenToken(),
				rightParenToken(),
				semiColonToken(),
				identToken("b"),
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
		{
			name: "CommaSeparatedNumbersParameters",
			code: Str("test(0X6,0x7,0x78,0X69,8,69,669,6.9,.9,3e1,4E7,4e7)"),
			want: tokens(
				identToken("test"),
				leftParenToken(),
				hexToken("0X6"),
				commaToken(),
				hexToken("0x7"),
				commaToken(),
				hexToken("0x78"),
				commaToken(),
				hexToken("0X69"),
				commaToken(),
				decimalToken("8"),
				commaToken(),
				decimalToken("69"),
				commaToken(),
				decimalToken("669"),
				commaToken(),
				decimalToken("6.9"),
				commaToken(),
				decimalToken(".9"),
				commaToken(),
				decimalToken("3e1"),
				commaToken(),
				decimalToken("4E7"),
				commaToken(),
				decimalToken("4e7"),
				rightParenToken(),
			),
		},
		{
			name: "CommaSeparatedNumbersAndStringsParameters",
			code: Str(`test("",5,"i",4,"k",6.6,0x5,"jssucks")`),
			want: tokens(
				identToken("test"),
				leftParenToken(),
				stringToken(""),
				commaToken(),
				decimalToken("5"),
				commaToken(),
				stringToken("i"),
				commaToken(),
				decimalToken("4"),
				commaToken(),
				stringToken("k"),
				commaToken(),
				decimalToken("6.6"),
				commaToken(),
				hexToken("0x5"),
				commaToken(),
				stringToken("jssucks"),
				rightParenToken(),
			),
		},
		{
			name: "PassingIdentifierAsArg",
			code: Str("test(arg)"),
			want: tokens(
				identToken("test"),
				leftParenToken(),
				identToken("arg"),
				rightParenToken(),
			),
		},
		{
			name: "PassingIdentifiersAsArg",
			code: Str("test(arg,arg2,i4k)"),
			want: tokens(
				identToken("test"),
				leftParenToken(),
				identToken("arg"),
				commaToken(),
				identToken("arg2"),
				commaToken(),
				identToken("i4k"),
				rightParenToken(),
			),
		},
		{
			name: "CommaSeparatedEverything",
			code: Str(`test("",5,"i",4,"k",6.6,0x5,arg,"jssucks")`),
			want: tokens(
				identToken("test"),
				leftParenToken(),
				stringToken(""),
				commaToken(),
				decimalToken("5"),
				commaToken(),
				stringToken("i"),
				commaToken(),
				decimalToken("4"),
				commaToken(),
				stringToken("k"),
				commaToken(),
				decimalToken("6.6"),
				commaToken(),
				hexToken("0x5"),
				commaToken(),
				identToken("arg"),
				commaToken(),
				stringToken("jssucks"),
				rightParenToken(),
			),
		},
	})
}

func TestPosition(t *testing.T) {
	cases := []TestCase{
		{
			name:          "MinusDecimal",
			code:          Str("-1"),
			checkPosition: true,
			want:          tokens(minusTokenPos(1, 1), decimalTokenPos("1", 1, 2)),
		},
		{
			name:          "PlusDecimal",
			code:          Str("+1"),
			checkPosition: true,
			want:          tokens(plusTokenPos(1, 1), decimalTokenPos("1", 1, 2)),
		},
		{
			name:          "PlusMinusDecimal",
			code:          Str("+-666"),
			checkPosition: true,
			want:          tokens(plusTokenPos(1, 1), minusTokenPos(1, 2), decimalTokenPos("666", 1, 3)),
		},
	}

	for _, lineTerminator := range lineTerminators() {
		lt := lineTerminator.val
		code := sfmt(`func(a)%sfuncb(1)%sfuncc("hi")`, lt, lt)
		cases = append(cases, TestCase{
			name:          "FuncallsSeparatedBy" + lineTerminator.name,
			code:          code,
			checkPosition: true,
			want: tokens(
				identTokenPos("func", 1, 1),
				leftParenTokenPos(1, 5),
				identTokenPos("a", 1, 6),
				rightParenTokenPos(1, 7),
				ltTokenPos(lt, 1, 8),
				identTokenPos("funcb", 2, 1),
				leftParenTokenPos(2, 6),
				decimalTokenPos("1", 2, 7),
				rightParenTokenPos(2, 8),
				ltTokenPos(lt, 2, 9),
				identTokenPos("funcc", 3, 1),
				leftParenTokenPos(3, 6),
				stringTokenPos("hi", 3, 7),
				rightParenTokenPos(3, 11),
			),
		})
	}

	runTests(t, cases)
}

func TestIllegalIdentifiers(t *testing.T) {
	t.Skip("TODO")
}

func TestIllegalSingleDot(t *testing.T) {
	cases := []TestCase{
		{
			name: "Nothing",
			code: Str("."),
			want: []lexer.Tokval{illegalToken(".")},
		},
	}

	for _, ts := range tokenSeparators() {
		code := sfmt(".%s.", ts.Value.String())
		cases = append(cases, TestCase{
			name: ts.Type.String(),
			code: code,
			want: []lexer.Tokval{illegalToken(code.String())},
		})
	}

	runTests(t, cases)
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

func TestInvalidStrings(t *testing.T) {

	cases := []TestCase{
		{
			name: "SingleDoubleQuote",
			code: Str(`"`),
			want: []lexer.Tokval{illegalToken(`"`)},
		},
		{
			name: "NoEndingDoubleQuote",
			code: Str(`"dsadasdsa123456`),
			want: []lexer.Tokval{illegalToken(`"dsadasdsa123456`)},
		},
	}

	for _, lineTerminator := range lineTerminators() {
		code := fmt.Sprintf(`"head%stail"`, lineTerminator.val)
		cases = append(cases, TestCase{
			code: Str(code),
			name: "NewlineTerminator" + lineTerminator.name,
			want: []lexer.Tokval{illegalToken(code)},
		})
	}

	runTests(t, cases)
}

func TestIllegalNumericLiterals(t *testing.T) {

	corruptedHex := messStr(Str("0x01234"), 4)
	corruptedDecimal := messStr(Str("1234"), 3)
	corruptedNumber := messStr(Str("0"), 1)

	cases := []TestCase{
		{
			name: "IncompleteExponentPart",
			code: Str("1e"),
			want: []lexer.Tokval{
				illegalToken("1e"),
			},
		},
		{
			name: "IncompleteUpperExponentPart",
			code: Str("1E"),
			want: []lexer.Tokval{
				illegalToken("1E"),
			},
		},
		{
			name: "IncompleteExponentPartByComma",
			code: Str("1e,"),
			want: []lexer.Tokval{
				illegalToken("1e,"),
			},
		},
		{
			name: "IncompleteExponentPartByRightParen",
			code: Str("1e)"),
			want: []lexer.Tokval{
				illegalToken("1e)"),
			},
		},
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
			name: "BigRealWithTwoDots",
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
	}

	for _, lineTerminator := range lineTerminators() {
		invalidReal := sfmt(".%s5", lineTerminator.val)
		invalidHexa := sfmt("0x%sFF", lineTerminator.val)
		invalidExp := sfmt("1e%s1", lineTerminator.val)

		newcases := []TestCase{
			{
				name: fmt.Sprintf("Invalid%sOnRealDecimal", lineTerminator.name),
				code: invalidReal,
				want: []lexer.Tokval{illegalToken(invalidReal.String())},
			},
			{
				name: fmt.Sprintf("Invalid%sOnHexaDecimal", lineTerminator.name),
				code: invalidHexa,
				want: []lexer.Tokval{illegalToken(invalidHexa.String())},
			},
			{
				name: fmt.Sprintf("Invalid%sOnExpDecimal", lineTerminator.name),
				code: invalidExp,
				want: []lexer.Tokval{illegalToken(invalidExp.String())},
			},
		}

		cases = append(cases, newcases...)
	}

	runTests(t, cases)
}

func TestNoOutputFor(t *testing.T) {
	runTests(t, []TestCase{
		{
			name: "EmptyString",
			code: Str(""),
			want: []lexer.Tokval{EOF},
		},
	})
}

func TestCorruptedInput(t *testing.T) {
	runTests(t, []TestCase{
		{
			name: "AtStart",
			code: messStr(Str(""), 0),
			want: []lexer.Tokval{illegalToken(messStr(Str(""), 0).String())},
		},
	})
}

type LineTerminator struct {
	name string
	val  string
}

func lineTerminators() []LineTerminator {
	return []LineTerminator{
		{name: "LineFeed", val: "\u000A"},
		{name: "CarriageReturn", val: "\u000D"},
		{name: "LineSeparator", val: "\u2028"},
		{name: "ParagraphSeparator", val: "\u2029"},
	}
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

// runTokenSepTests will take an array of test cases and run the
// tests with different token separators intertwined between
// the wanted tokens of each test case, validating if the tokens
// gets separated correctly.
//
// This functions is useful to reuse tokens test cases to validate
// token separation/splitting (with newlines or semicolons for example).
func runTokenSepTests(t *testing.T, testcases []TestCase) {
	for _, ts := range tokenSeparators() {
		runTests(t, intertwineOnTestCases(ts, testcases))
	}
}

func tokenSeparators() []lexer.Tokval {
	tokens := []lexer.Tokval{}
	for _, lt := range lineTerminators() {
		tokens = append(tokens, ltToken(lt.val))
	}
	tokens = append(tokens, semiColonToken())
	tokens = append(tokens, rightParenToken())
	tokens = append(tokens, commaToken())
	return tokens
}

func illegalToken(val string) lexer.Tokval {
	return lexer.Tokval{
		Type:  token.Illegal,
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
			t.Errorf("\nwanted:\ntoken[%d][%v]\n\ngot:\ntoken[%d][%v]", i, w, i, g)
			t.Errorf("\nwanted:\n%v\ngot:\n%v\n", tc.want, got)
		}

		if tc.checkPosition {
			if !w.EqualPos(g) {
				t.Errorf("\nwant=%+v\ngot=%+v\nare equal but dont have the same position", w, g)
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
// provided on tcases, generating a new array of TestCases.
//
// The array of TestCases is generated by prepending code and the
// wanted tokens from the given tcase on each test case on tcases.
// EOF should not be provided on the
// given tcase since it will be prepended on each test case inside given tcases
// and the provided tcases will already have their own EOF.
func prependOnTestCases(tcase TestCase, tcases []TestCase) []TestCase {
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

// intertwineOnTestCases will intertwine the given token on each test case
// provided by tcases. For each test case it will take all the wanted tokens
// and interwine the provided token between them. The resulting test case
// will have the token that has been interwined on the expectation also.
// If the test case has only one expected token on its wanted list the token
// will be duplicated so it can be intertwined.
//
// This functions is useful to test easily the handling of tokens that
// acts as generic separators between other tokens, like semi colons/spaces/newlines.
//
// All information regarding token positions is ignored.
func intertwineOnTestCases(tok lexer.Tokval, tcases []TestCase) []TestCase {
	newCases := make([]TestCase, len(tcases))

	for i, tcase := range tcases {
		name := fmt.Sprintf("%s/IntertwinedWith%s", tcase.name, tok.Type)
		want, hasEOF := removeEOF(tcase.want)

		if len(want) == 1 {
			want = append(want, want[0])
		}

		newwant := []lexer.Tokval{}

		for i := 0; i < len(want)-1; i++ {
			newwant = append(newwant, want[i])
			newwant = append(newwant, tok)
		}

		newwant = append(newwant, want[len(want)-1])
		if hasEOF {
			newwant = append(newwant, EOF)
		}

		newCases[i] = newTestCase(name, newwant)
	}
	return newCases
}

func newTestCase(name string, tokens []lexer.Tokval) TestCase {
	tcase := TestCase{
		name: name,
		want: tokens,
	}

	par := Str(`"`)

	for _, tok := range tokens {
		if tok.Equal(EOF) {
			continue
		}

		if tok.Type != token.String {
			tcase.code = tcase.code.Append(tok.Value)
		} else {
			// WHY: on strings the parenthesis is removed when the token is produced
			tcase.code = tcase.code.Append(par)
			tcase.code = tcase.code.Append(tok.Value)
			tcase.code = tcase.code.Append(par)
		}
	}

	return tcase
}

func removeEOF(tokens []lexer.Tokval) ([]lexer.Tokval, bool) {
	if len(tokens) == 0 {
		return tokens, false
	}

	lasttoken := tokens[len(tokens)-1]
	if lasttoken.Equal(EOF) {
		// WHY: got nasty side effects bugs if dont copy tokens array here
		// the provided slice underlying array is modified and all hell break loses =D
		newtokens := make([]lexer.Tokval, len(tokens)-1)
		copy(newtokens, tokens)
		return newtokens, true
	}

	return tokens, false
}

func sfmt(format string, a ...interface{}) utf16.Str {
	return Str(fmt.Sprintf(format, a...))
}

func tokval(t token.Type, val string, line uint, column uint) lexer.Tokval {
	return lexer.Tokval{
		Type: t,
		Value: Str(val),
		Line: line,
		Column: column,
	}
}

func nullToken() lexer.Tokval {
	return tokval(token.Null, "null", 0, 0)
}

func minusToken() lexer.Tokval {
	return lexer.Tokval{
		Type:  token.Minus,
		Value: Str("-"),
	}
}

func plusToken() lexer.Tokval {
	return lexer.Tokval{
		Type:  token.Plus,
		Value: Str("+"),
	}
}

func leftParenToken() lexer.Tokval {
	return leftParenTokenPos(0, 0)
}

func leftParenTokenPos(line uint, column uint) lexer.Tokval {
	return lexer.Tokval{
		Type:   token.LParen,
		Value:  Str("("),
		Line:   line,
		Column: column,
	}
}

func rightParenToken() lexer.Tokval {
	return rightParenTokenPos(0, 0)
}

func rightParenTokenPos(line uint, column uint) lexer.Tokval {
	return lexer.Tokval{
		Type:   token.RParen,
		Value:  Str(")"),
		Line:   line,
		Column: column,
	}
}

func minusTokenPos(line uint, column uint) lexer.Tokval {
	return lexer.Tokval{
		Type:   token.Minus,
		Value:  Str("-"),
		Line:   line,
		Column: column,
	}
}

func plusTokenPos(line uint, column uint) lexer.Tokval {
	return lexer.Tokval{
		Type:   token.Plus,
		Value:  Str("+"),
		Line:   line,
		Column: column,
	}
}

func decimalTokenPos(dec string, line uint, column uint) lexer.Tokval {
	return lexer.Tokval{
		Type:   token.Decimal,
		Value:  Str(dec),
		Line:   line,
		Column: column,
	}
}

func decimalToken(dec string) lexer.Tokval {
	return lexer.Tokval{
		Type:  token.Decimal,
		Value: Str(dec),
	}
}

func dotToken() lexer.Tokval {
	return lexer.Tokval{
		Type:  token.Dot,
		Value: Str("."),
	}
}

func hexToken(hex string) lexer.Tokval {
	return lexer.Tokval{
		Type:  token.Hexadecimal,
		Value: Str(hex),
	}
}

func semiColonToken() lexer.Tokval {
	return semicolonTokenPos(0, 0)
}

func semicolonTokenPos(line uint, column uint) lexer.Tokval {
	return lexer.Tokval{
		Type:   token.SemiColon,
		Value:  Str(";"),
		Line:   line,
		Column: column,
	}
}

func stringToken(s string) lexer.Tokval {
	return stringTokenPos(s, 0, 0)
}

func stringTokenPos(s string, line uint, column uint) lexer.Tokval {
	return lexer.Tokval{
		Type:   token.String,
		Value:  Str(s),
		Line:   line,
		Column: column,
	}
}

func identToken(s string) lexer.Tokval {
	return identTokenPos(s, 0, 0)
}

func identTokenPos(s string, line uint, column uint) lexer.Tokval {
	return lexer.Tokval{
		Type:   token.Ident,
		Value:  Str(s),
		Line:   line,
		Column: column,
	}
}

func ltToken(s string) lexer.Tokval {
	return ltTokenPos(s, 0, 0)
}

func ltTokenPos(s string, line uint, column uint) lexer.Tokval {
	return lexer.Tokval{
		Type:   token.Newline,
		Value:  Str(s),
		Line:   line,
		Column: column,
	}
}

func commaToken() lexer.Tokval {
	return lexer.Tokval{
		Type:  token.Comma,
		Value: Str(","),
	}
}

func tokens(t ...lexer.Tokval) []lexer.Tokval {
	return append(t, EOF)
}
