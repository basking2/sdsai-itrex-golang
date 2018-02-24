package itrex

import (
	"github.com/basking2/sdsai-itrex-golang/pkg/sdsai/iterator"
)

type FunctionInterface func(iterator.Iterator, *Context) interface{}
