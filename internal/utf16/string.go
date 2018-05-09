package utf16

type (
	// Str is a UTF-16 encoded string
	Str []uint16
)

// String is the UTF-8 string representation of Str
func (s Str) String() string {
	return Decode(s)
}