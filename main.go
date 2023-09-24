package main

import (
	"bytes"
	"cmp"
	"flag"
	"fmt"
	"go/ast"
	"go/types"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/donatorsky/go-cmder/internal/template"
	"github.com/donatorsky/go-cmder/internal/utils"
	"golang.org/x/tools/go/packages"
)

var filenamePattern = regexp.MustCompile(`(ID|JSON|URL|[[:upper:]])`)

func main() {
	logger := slog.NewLogLogger(slog.NewTextHandler(os.Stdout, nil), slog.LevelError)

	params := newParams()

	flag.BoolVar(&params.mutable, "mutable", false, "Whether the generated command should be mutable.")
	flag.BoolVar(&params.includeUnexported, "include-unexported", false, "Whether to include unexported fields.")
	flag.BoolVar(&params.sorted, "sorted", false, "Whether to generate fields in alphabetic ascending order.")
	flag.StringVar(&params.out, "out", "", "Where write to the generated command.")
	flag.Var(params.exclude, "exclude", "Struct field's name to ignore when generating command.")
	flag.Var(params.include, "include", "Struct field's name to generate command from. Overrides -exclude flag.")
	flag.Var(params.constructor, "constructor", `Constructor name and comma-separated list of fields.
Use "default" as a constructor name to generate default constructor.

E.g.:
-constructor default:foo,bar CreateStructCmd       // Generates NewCreateStructCmd(foo fooType, bar barType)
-constructor WithFooAndBar:foo,bar CreateStructCmd // Generates NewCreateStructCmdWithFooAndBar(foo fooType, bar barType)`)

	flag.Usage = func() {
		fmt.Println(`go-cmder [flags] struct CommandName`)
		flag.PrintDefaults()
	}

	flag.Parse()

	if flag.NArg() < 2 {
		logger.Fatalln("Missing required arguments")
	}

	params.structName = flag.Arg(0)
	params.commandName = flag.Arg(1)

	if params.out == "" {
		params.out = fmt.Sprintf(
			"%s.go",
			strings.ToLower(strings.Trim(filenamePattern.ReplaceAllString(params.commandName, "_$1"), "_")),
		)
	}

	cwd, err := os.Getwd()
	if err != nil {
		logger.Fatalf("Could not get current working directory: %s\n", err)
	}

	params.out = filepath.Join(cwd, params.out)

	pkgs, err := packages.Load(&packages.Config{
		Mode: packages.NeedName | packages.NeedSyntax | packages.NeedTypes,
		Dir:  cwd,
	})
	if err != nil {
		logger.Fatalf("Could not load source package: %s\n", err)
	}
	if len(pkgs) == 0 {
		logger.Fatalln("package not found")
	}
	if len(pkgs) > 1 {
		logger.Fatalln("found more than one package")
	}
	if len(pkgs[0].Errors) > 0 {
		logger.Fatalf("failures: %s\n", pkgs[0].Errors)
	}

	importsAliases := getImportsAliases(pkgs[0].Syntax)

	obj := pkgs[0].Types.Scope().Lookup(params.structName)
	if obj == nil {
		logger.Fatalf("struct %s not found\n", params.structName)
	}

	structType, ok := obj.Type().Underlying().(*types.Struct)
	if !ok {
		logger.Fatalf("%s (%s) is not a struct", params.structName, obj.Type())
	}

	fields := utils.NewUniqueSlice[*template.FieldData](
		utils.UniqueSliceWithCapacity(uint(structType.NumFields() - params.exclude.Len())),
	)

	imports := utils.NewUniqueSlice[template.CommandDataImport]()

	for i := 0; i < structType.NumFields(); i++ {
		field := structType.Field(i)

		if params.include.Empty() {
			if params.exclude.Has(field.Name()) {
				continue
			}
		} else {
			if !params.include.Has(field.Name()) {
				continue
			}
		}

		if !params.includeUnexported && !field.Exported() {
			continue
		}

		var commandDataField *template.FieldData
		if namedType, ok := field.Type().(*types.Named); ok {
			typeFQN := namedType.String()
			typeNameIndex := strings.LastIndexByte(typeFQN, '.')
			if typeNameIndex == -1 {
				logger.Fatalf("Type FQN %q is expected to contain .", typeFQN)
			}

			var typeNamespace string
			packageName := typeFQN[:typeNameIndex]

			if alias, exists := importsAliases[packageName]; exists {
				typeNamespace = alias
				_, _ = imports.Append(template.CommandDataImport{
					Alias:   &alias,
					Package: packageName,
				})
			} else {
				typeNamespace = packageName[strings.LastIndexByte(packageName, '/')+1:]
				_, _ = imports.Append(template.CommandDataImport{
					Alias:   nil,
					Package: packageName,
				})
			}

			commandDataField = &template.FieldData{
				CommandName: params.commandName,
				Mutable:     params.mutable,
				Name:        field.Name(),
				Type:        fmt.Sprintf("%s.%s", typeNamespace, typeFQN[typeNameIndex+1:]),
			}
		} else {
			commandDataField = &template.FieldData{
				CommandName: params.commandName,
				Mutable:     params.mutable,
				Name:        field.Name(),
				Type:        field.Type().String(),
			}
		}

		if fields.Has(commandDataField) {
			logger.Fatalf("Fields' names conflict with %q.", commandDataField.Name)
		}

		_, _ = fields.Append(commandDataField)
	}

	if params.sorted {
		fields.Sort(func(i, j *template.FieldData) int {
			return cmp.Compare(i.Name, j.Name)
		})
	}

	tpl, err := template.NewTemplate()
	if err != nil {
		logger.Fatalf("Failed to parse command template: %s\n", err)
	}

	var b bytes.Buffer

	var constructors []string
	for _, constructor := range params.constructor.Items() {
		constructorData := template.ConstructorData{
			CommandName: params.commandName,
			Mutable:     params.mutable,
			Name:        params.commandName,
			Fields:      make([]*template.FieldData, 0, len(constructor.Params)),
		}

		if strings.ToLower(constructor.Name) != "default" {
			constructorData.Name = fmt.Sprintf("%s%s", params.commandName, constructor.Name)
		}

		for _, param := range constructor.Params {
			fieldData := &template.FieldData{
				Name: param,
			}

			if err := fields.GetByItem(&fieldData); err != nil {
				logger.Fatalf("Cannot build %s constructor: field %s does not exist, is excluded or not included", constructor.Name, param)
			}

			constructorData.Fields = append(constructorData.Fields, fieldData)
		}

		if err := tpl.ExecuteConstructorTemplate(&b, &constructorData); err != nil {
			logger.Fatalf("Failed to generate command getter: %s\n", err)
		}

		constructors = append(constructors, b.String())
		b.Reset()
	}

	var methods []string

	for _, field := range fields.Items() {
		if err := tpl.ExecuteGetterTemplate(&b, field); err != nil {
			logger.Fatalf("Failed to generate command getter: %s\n", err)
		}

		methods = append(methods, b.String())
		b.Reset()

		if err := tpl.ExecuteSetterTemplate(&b, field); err != nil {
			logger.Fatalf("Failed to generate command setter: %s\n", err)
		}

		methods = append(methods, b.String())
		b.Reset()

		if err := tpl.ExecuteHaserTemplate(&b, field); err != nil {
			logger.Fatalf("Failed to generate command haser: %s\n", err)
		}

		methods = append(methods, b.String())
		b.Reset()
	}

	file, err := os.OpenFile(params.out, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		logger.Fatalf("Failed to open command output file for writing: %s\n", err)
	}

	if err := tpl.ExecuteCommandTemplate(file, &template.CommandData{
		PackageName:  pkgs[0].Name,
		Imports:      imports.Items(),
		CommandName:  params.commandName,
		Fields:       fields.Items(),
		Constructors: constructors,
		Methods:      methods,
	}); err != nil {
		logger.Fatalf("Failed to write command output file: %s\n", err)
	}
}

func newParams() params {
	return params{
		exclude: utils.NewUniqueMultiFlag(utils.StringSetter),
		include: utils.NewUniqueMultiFlag(utils.StringSetter),
		constructor: utils.NewUniqueMultiFlag(
			func(value string) (c constructor, _ error) {
				nameAndParams := strings.SplitN(value, ":", 2)

				c.Name = nameAndParams[0]

				if len(nameAndParams) == 2 {
					c.Params = strings.Split(nameAndParams[1], ",")
				}

				return
			},
			utils.UniqueSliceWithOnDuplicateKeyError(func(key, item any) error {
				return fmt.Errorf("duplicated constructor name %q", key)
			}),
		),
	}
}

type params struct {
	mutable           bool
	includeUnexported bool
	sorted            bool
	out               string
	exclude           *utils.UniqueMultiFlag[string]
	include           *utils.UniqueMultiFlag[string]
	constructor       *utils.UniqueMultiFlag[constructor]
	structName        string
	commandName       string
}

type constructor struct {
	Name   string
	Params []string
}

func (c constructor) UniqueValue() any {
	return c.Name
}

func getImportsAliases(syntaxTree []*ast.File) map[string]string {
	aliases := make(map[string]string)

	for _, syntax := range syntaxTree {
		for _, importSpec := range syntax.Imports {
			if importSpec.Name != nil && importSpec.Name.Name != "." && importSpec.Name.Name != "_" {
				aliases[strings.Trim(importSpec.Path.Value, `"`)] = importSpec.Name.Name
			}
		}
	}

	return aliases
}
