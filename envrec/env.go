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
		New(name utf16.Str, candelete bool) error
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

func (env *Decl) New(name utf16.Str, candelete bool) error {
	if len(name) == 0 {
		return fmt.Errorf("empty binding name")
	}

	env.records[name.String()] = Record{
		mutable:   true,
		deletable: candelete,
		value:     types.Undef,
	}

	return nil
}

func (env *Decl) Has(name utf16.Str) bool {
	_, ok := env.records[name.String()]
	return ok
}

func (env *Decl) Set(name utf16.Str, v types.Value, musterr bool) error {
	if !env.Has(name) {
		if musterr {
			return fmt.Errorf("%s is not defined", name)
		}

		env.New(name, true)
	}

	str := name.String()
	r := env.records[str]
	r.value = v
	env.records[str] = r
	return nil
}

func (env *Decl) Get(name utf16.Str, musterr bool) (types.Value, error) {
	r, ok := env.records[name.String()]
	if !ok {
		if musterr {
			return nil, fmt.Errorf("%s is not defined", name)
		}

		return types.Undef, nil
	}

	return r.value, nil
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
