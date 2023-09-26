package types

import (
	"fmt"
	"go/types"
	"strings"

	"github.com/donatorsky/go-cmder/internal/utils"
	"golang.org/x/tools/go/packages"
)

func NewRegistry(pkg *packages.Package) *Registry {
	r := Registry{
		selfPkg: pkg.PkgPath,
		types:   map[string]*Type{},
		imports: utils.NewUniqueSlice[*Type](),
	}

	for _, p := range pkg.Imports {
		r.types[p.PkgPath] = &Type{
			Alias: nil,
			Name:  p.Name,
			Path:  p.PkgPath,
		}
	}

	for _, syntax := range pkg.Syntax {
		for _, importSpec := range syntax.Imports {
			if importSpec.Name == nil || importSpec.Name.Name == "." || importSpec.Name.Name == "_" {
				continue
			}

			t, ok := r.types[strings.Trim(importSpec.Path.Value, `"`)]
			if !ok {
				continue
			}

			t.Alias = &importSpec.Name.Name
		}
	}

	return &r
}

type Registry struct {
	selfPkg string
	types   map[string]*Type
	imports *utils.UniqueSlice[*Type]
}

func (r *Registry) Imports() []*Type {
	return r.imports.Items()
}

func (r *Registry) Resolve(fieldType types.Type) (pointer string, unwrappedType string, _ error) {
	for {
		pointerType, ok := fieldType.(*types.Pointer)
		if !ok {
			break
		}

		pointer += "*"
		fieldType = pointerType.Elem()
	}

	switch actualType := fieldType.(type) {
	case *types.Named,
		*types.Struct,
		*types.Signature:
		typeFQN := actualType.String()

		for name, t := range r.types {
			if !strings.Contains(typeFQN, fmt.Sprintf("%s.", name)) {
				continue
			}

			alias := t.Name
			if t.Alias != nil {
				alias = *t.Alias
			}

			typeFQN = strings.ReplaceAll(typeFQN, fmt.Sprintf("%s.", name), fmt.Sprintf("%s.", alias))

			_, _ = r.imports.Append(t)
		}

		return pointer, strings.ReplaceAll(typeFQN, fmt.Sprintf("%s.", r.selfPkg), ""), nil

	case *types.Slice:
		elemPointer, elemType, err := r.Resolve(actualType.Elem())
		if err != nil {
			return "", "", err
		}

		return pointer, fmt.Sprintf("[]%s%s", elemPointer, elemType), nil

	case *types.Array:
		elemPointer, elemType, err := r.Resolve(actualType.Elem())
		if err != nil {
			return "", "", err
		}

		return pointer, fmt.Sprintf("[%d]%s%s", actualType.Len(), elemPointer, elemType), nil

	default:
		return pointer, fieldType.String(), nil
	}
}
