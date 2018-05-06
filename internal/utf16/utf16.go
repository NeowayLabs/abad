package utf16

import "unicode/utf16"

// Encode an UTF-8 string into an UTF-16 Str.
func Encode(a string) Str {
	return Str(utf16.Encode([]rune(a)))
}

// Decode an UTF-16 string into a UTF-8 string
func Decode(a Str) string {
	return string(utf16.Decode([]uint16(a)))
}