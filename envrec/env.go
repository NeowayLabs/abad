package envrec

import (
	"fmt"

	"github.com/NeowayLabs/abad/internal/utf16"
	"github.com/NeowayLabs/abad/types"
)

type (
	// Env is an environment record storage
	Env interface {
		Has(name utf16.Str) bool
		New(name utf16.Str, candelete bool)
		Set(name utf16.Str, value types.Value, musterr bool) error
		Get(name utf16.Str, musterr bool) (types.Value, error)
		Del(name utf16.Str) bool

		ImplicitThis() types.Value
	}

	Record struct {
		mutable   bool
		deletable bool
		value     types.Value
	}

	// Declarative environment record
	// https://es5.github.io/#x10.2.1
	Decl struct {
		records map[string]Record
	}
)

func NewDeclEnv() *Decl {
	return &Decl{
		records: make(map[string]Record),
	}
}

func (env *Decl) New(name utf16.Str, candelete bool) {
	env.records[name.String] = Record{
		mutable:   true,
		deletable: candelete,
		value:     types.Undef,
	}
}

func (env *Decl) Has(name utf16.Str) bool {
	_, ok := env.records[name.String()]
	return ok
}

func (env *Decl) Set(name utf16.Str, v types.Value, musterr bool) error {
	if !env.Has(name) {
		if musterr {
			return fmt.Errorf("%s is not defined", n)
		}

		env.New(n, true)
	}

	str := n.String()
	r := env.records[str]
	r.Value = v
	env.records[str] = r
	return nil
}

func (env *Decl) Get(name utf16.Str, musterr bool) (types.Value, error) {
	r, ok := env.records[name.String()]
	if !ok {
		if musterr {
			return nil, fmt.Errorf("%s is not defined", n)
		}

		return types.Undef
	}

	return r.Value, nil
}

func (env *Decl) Del(name utf16.Str) bool {
	if !env.Has(name) {
		return false
	}

	r := env.records[name.String()]
	if !r.deletable {
		return false
	}

	delete(env.records, name.String())
	return true
}

func (env *Decl) ImplicitThis() types.Value {
	return nil
}
