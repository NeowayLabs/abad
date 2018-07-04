package types

import (
	"fmt"
)

type (
	TypeError struct {
		msg string
	}
)

func NewTypeError(format string, args ...interface{}) TypeError {
	err := TypeError{
		msg: fmt.Sprintf(format, args...),
	}

	return err
}

func (e TypeError) Error() string {
	// TODO(i4k): improve this
	return fmt.Sprintf("TypeError: %s\n\tat anonymous:1:1", e.msg)
}

func (e TypeError) Exception() bool { return true }
