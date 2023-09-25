package examples

import (
	"time"

	baz "github.com/donatorsky/go-cmder/examples/bar"
	"github.com/donatorsky/go-cmder/examples/foo"
	"github.com/gofrs/uuid/v5"
	"gopkg.in/yaml.v3"
)

//go:generate go-cmder -out create_struct_cmd.go -exclude Ignore -sorted -constructor default -constructor WithIntAndString:Int,String Struct CreateStructCmd
//go:generate go-cmder -out struct_with_includes_cmd.go -exclude Ignore -include Ignore -include Int -include String -include UniqueMultiFlag -sorted -constructor default -constructor WithIntAndStringAndUniqueMultiFlag:Int,String,UniqueMultiFlag Struct StructWithIncludesCmd
//go:generate go-cmder -out mutable_update_struct_cmd.go -sorted -constructor default -constructor WithIntAndUniqueMultiFlag:Int,UniqueMultiFlag -mutable Struct UpdateStructCmd
//go:generate go-cmder -sorted -constructor default Struct StructWithDefaultOutputFilenameCmd
type Struct struct {
	String       string
	StringPtr    *string
	StringPtrPtr **string
	string       string
	Ignore       string
	Int          int
	IntPtr       *int
	Float        float64
	FloatPtr     *float64
	Slice        []string
	SliceOfPtrs  []*string
	SlicePtr     *[]string
	Array        [3]string
	ArrayOfPtrs  [3]*string
	ArrayPtr     *[3]string
	Map          map[string]int
	MapOfPtrs    *map[string]*int
	MapPtr       *map[string]int
	Time         time.Time
	TimePtr      *time.Time
	Struct       struct {
		Foo string
		Bar int
	}
	StructPtr *struct {
		Foo string
		Bar int
	}
	FieldData          foo.OtherStruct
	FieldDataPtr       *foo.OtherStruct
	UniqueMultiFlag    baz.OtherGenericStruct[any]
	UniqueMultiFlagPtr *baz.OtherGenericStruct[any]
	Uuid               uuid.UUID
	UuidPtr            *uuid.UUID
	Yaml               yaml.Decoder
	YamlPtr            *yaml.Decoder
}
