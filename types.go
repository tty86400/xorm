package xorm

import (
	"reflect"

	"github.com/chinahdkj/core"
)

var (
	ptrPkType = reflect.TypeOf(&core.PK{})
	pkType    = reflect.TypeOf(core.PK{})
)
