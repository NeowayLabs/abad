package utf16

type (
	// Str is a UTF-16 encoded string
	Str []uint16
)

func S(a string) Str {
	return NewStr(a)
}

func NewStr(a string) Str {
	return Encode(a)
}

// String is the UTF-8 string representation of Str
func (s Str) String() string {
	return Decode(s)
}

// Contains checks if substr is inside s
func (s Str) Contains(substr Str) bool {
	return s.Index(substr) >= 0 
}

// Index the position of substr in s
func (s Str) Index(substr Str) int {
	switch n := len(substr); {
	case n == 0:
		return 0
	case n == 1:
		for i := range s {
			if s[i] == substr[0] {
				return i
			}
		}
		return -1
	case n == len(s):
		if s.Equal(substr) {
			return 0
		}
		return -1
	case n > len(s):
		return -1
	}

	// brute force
	var i, j int
	for j < len(substr) && i < len(s) {
		si := s[i]
		bi := substr[j]

		i++

		if si == bi {
			j++
		} else {
			j = 0
		}
	}

	if j == len(substr) {
		return i - len(substr)
	}
	return -1

}

// Equal checks if s is equal o
func (s Str) Equal(o Str) bool {
	if len(s) != len(o) {
		return false
	}

	for i := 0; i < len(s); i++ {
		if s[i] != o[i] {
			return false
		}
	}
	return true
}

func (s Str) TrimPrefix(substr Str) Str {
	if s.Index(substr) == 0 {
		return Str(s[len(substr):])
	}
	return s
}