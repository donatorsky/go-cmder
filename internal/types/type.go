package types

import "fmt"

type Type struct {
	Alias *string
	Name  string
	Path  string
}

func (t *Type) UniqueValue() any {
	return fmt.Sprintf("%#v", t)
}
