package examples

import (
	"time"

	"github.com/donatorsky/go-cmder/internal/template"
	myAlias "github.com/donatorsky/go-cmder/internal/utils"
)

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
