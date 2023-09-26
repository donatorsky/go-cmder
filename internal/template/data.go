package template

import (
	"strings"

	internalTypes "github.com/donatorsky/go-cmder/internal/types"
)

type CommandData struct {
	PackageName  string
	Imports      []*internalTypes.Type
	CommandName  string
	Fields       []*FieldData
	Constructors []string
	Methods      []string
}

type FieldData struct {
	CommandName string
	Mutable     bool
	Name        string
	Pointer     string
	Type        string
}

func (c *FieldData) UniqueValue() any {
	return strings.ToLower(c.Name)
}

type ConstructorData struct {
	CommandName string
	Mutable     bool
	Name        string
	Fields      []*FieldData
}

func (c *ConstructorData) UniqueValue() any {
	return strings.ToLower(c.Name)
}
