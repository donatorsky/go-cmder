package examples

import (
	"io/fs"
	"time"
)

type InternalStruct struct {
	Foo  string
	Time time.Time
	File fs.File
}
