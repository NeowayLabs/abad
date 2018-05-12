package types

import "github.com/NeowayLabs/abad/internal/utf16"

type (
	String utf16.Str
)

func NewString(str string) String {
	return String(utf16.Encode(str))
}