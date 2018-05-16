package types

import (
	"fmt"

	"github.com/NeowayLabs/abad/internal/utf16"
)

type (
	TypeError struct {
		*Object
	}
)

var messageAttr = utf16.S("message")

func NewTypeError(msg Value) *TypeError {
	err := &TypeError{
		Object: NewRawObject(),
	}

	err.DefineOwnProperty(messageAttr, NewDataPropDesc(
		msg, true, false, true,
	).ToObject(), false)

	return err
}

func NewTypeErrorS(format string, args ...interface{}) *TypeError {
	return NewTypeError(NewString(fmt.Sprintf(format, args...)))
}

func (e *TypeError) Error() string {
	msg, err := e.Get(messageAttr)
	if err != nil {
		msg = NewString("")
	}

	// TODO(i4k): improve this
	return fmt.Sprintf("TypeError: %s\n\tat anonymous:1:1",
		msg.ToString())
}

func (e TypeError) Exception() bool { return true }