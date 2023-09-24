# go-cmder

A command from struct generator.

## Installation

```shell
go install github.com/donatorsky/go-cmder@latest
```

## Usage

```shell
go-cmder [flags] struct CommandName
```

### Flags

#### `-constructor=name[:field1,fieldn...]`

Defines a name and comma-separated list of fields for command constructor.
Multiple usage allowed to generate many constructors.

A name is appended to the command name. E.g. given command name MyCmd, then `-constructor=WithFoo` generates `func NewMyCmdWithFoo`.
Use name `default` to generate constructor based on command name. E.g. given command name MyCmd, then `-constructor=default` generates `func NewMyCmd`.

List of fields can be omitted to generate constructor without fields.

#### `-exclude=field`

Excludes given struct field from command generation.
Multiple usage allowed.

#### `-include=field`

Includes given struct field in command generation.
When present, command is generated only from included fields.
Supersedes `-exclude` flag.
Multiple usage allowed.

#### `-include-unexported`

Includes unexported fields when generating a command.

#### `-mutable`

Generates a mutable command.
By default, commands are immutable, meaning that calling a setter returns new command instance.

#### `-out=path/to/file.go`

Generates a command in given file.
By default, command is generated in `CommandName.go` file.

#### `-sorted`

Sort fields by name when generating a command.

## Example
```go
package foobar

//go:generate go-cmder -out=create_struct_cmd.go -constructor=default Struct CreateStructCmd
type Struct struct {
	Foo string
	Bar int
}
```

Generates `create_struct_cmd.go` file:
```go
package foobar

type CreateStructCmd struct {
	vFoo string
	hasFoo bool

	vBar int
	hasBar bool
}

func NewCreateStructCmd() CreateStructCmd {
	return CreateStructCmd{}
}

func (cmd CreateStructCmd) Foo() string {
	return cmd.vFoo
}

func (cmd CreateStructCmd) SetFoo(v string) CreateStructCmd {
	cmd.hasFoo = true
	cmd.vFoo = v

	return cmd
}

func (cmd CreateStructCmd) HasFoo() bool {
	return cmd.hasFoo
}

func (cmd CreateStructCmd) Bar() int {
	return cmd.vBar
}

func (cmd CreateStructCmd) SetBar(v int) CreateStructCmd {
	cmd.hasBar = true
	cmd.vBar = v

	return cmd
}

func (cmd CreateStructCmd) HasBar() bool {
	return cmd.hasBar
}
```
