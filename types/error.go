package types

import (
	"fmt"

	"github.com/NeowayLabs/abad/internal/utf16"
)

type (
	TypeError struct {
		msg string
	}
)

var messageAttr = utf16.S("message")

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