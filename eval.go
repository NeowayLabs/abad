package abad

type (
	// TODO(i4k): Define Object interface
	Obj interface{}

	// Abad interpreter
	Abad struct{}
)

func (a *Abad) Exec(code string) (Obj, error) {
	return nil, nil
}
