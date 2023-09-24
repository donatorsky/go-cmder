package examples

import (
	"time"

	"github.com/donatorsky/go-cmder/internal/template"
	myAlias "github.com/donatorsky/go-cmder/internal/utils"
)

//go:generate go-cmder -out create_struct_cmd.go -exclude Ignore -sorted -constructor default -constructor WithIntAndString:Int,String Struct CreateStructCmd
//go:generate go-cmder -out struct_with_includes_cmd.go -exclude Ignore -include Ignore -include Int -include String -include UniqueMultiFlag -sorted -constructor default -constructor WithIntAndStringAndUniqueMultiFlag:Int,String,UniqueMultiFlag Struct StructWithIncludesCmd
//go:generate go-cmder -out mutable_update_struct_cmd.go -sorted -constructor default -constructor WithIntAndUniqueMultiFlag:Int,UniqueMultiFlag -mutable Struct UpdateStructCmd
//go:generate go-cmder -sorted -constructor default Struct StructWithDefaultOutputFilenameCmd
type Struct struct {
	String  string
	string  string
	Ignore  string
	Int     int
	Float   float64
	Slice   []string
	Array   [3]string
	Map     map[string]int
	Time    time.Time
	TimePtr *time.Time
	Struct  struct {
		Foo string
		Bar int
	}
	FieldData       template.FieldData
	UniqueMultiFlag myAlias.UniqueMultiFlag[any]
}
