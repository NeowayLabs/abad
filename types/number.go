package types

import "strconv"

type (
	Number float64
)

func (a Number) String() string {
	return strconv.FormatFloat(float64(a), 'e', -1, 64)
}