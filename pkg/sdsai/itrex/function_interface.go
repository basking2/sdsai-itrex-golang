package itrex

import (
	"github.com/basking2/sdsai-itrex-golang/pkg/sdsai/iterator"
)

type FunctionInterface interface {
	Apply(iterator.Iterator, *Context) interface{}
}
