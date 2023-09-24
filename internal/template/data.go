package template

import (
	"strings"
)

type CommandData struct {
	PackageName  string
	Imports      []CommandDataImport
	CommandName  string
	Fields       []*FieldData
	Constructors []string
	Methods      []string
}

type CommandDataImport struct {
	Alias   *string
	Package string
}

func (c *CommandDataImport) UniqueValue() any {
	return c.Package
}

type FieldData struct {
	CommandName string
	Mutable     bool
	Name        string
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
