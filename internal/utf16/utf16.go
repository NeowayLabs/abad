package utf16

import "unicode/utf16"

// Encode an UTF-8 string into an UTF-16 Str.
func Encode(a string) Str {
	return EncodeRunes([]rune(a))
}

// Encode runes into an UTF-16 Str
func EncodeRunes(r []rune) Str {
	return Str(utf16.Encode(r))
}

// Decode an UTF-16 string into a UTF-8 string
func Decode(a Str) string {
	return string(DecodeRunes(a))
}

func DecodeRunes(s Str) []rune {
	return utf16.Decode([]uint16(s))
}
