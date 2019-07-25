package xorm

import "github.com/chinahdkj/core"

type MutableFilter interface {
	AddFilter(filters ...core.Filter)
}
