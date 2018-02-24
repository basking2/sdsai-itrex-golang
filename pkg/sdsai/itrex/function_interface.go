package itrex

import (
	"github.com/basking2/sdsai-itrex-golang/pkg/sdsai/iterator"
)

type FunctionInterface interface {
	Apply(iterator.Iterator, *Context) interface{}
}

// A type to bind a Go function to an ItrEx function.
type BoundFunction struct {
	// The function definition.
	apply func(iterator.Iterator, *Context, interface{}) interface{}

	// Function state. This may be nil.
	ctx interface{}
}

// How to apply the function stored in a BoundFunction.
func (bf BoundFunction) Apply(i iterator.Iterator, ctx *Context) interface{} {
	return bf.apply(i, ctx, bf.ctx)
}

// Make a new FunctionInterface from a free function.
//
// apply - the function to bind.
// ctx - function state that is passed to the appy function.
//
// For functions with no state, the ctx may be nil, it is never read by
// this framework. It is simply passed to calls to the apply function as
// as way to close state.
func NewBoundFunction(apply func(iterator.Iterator, *Context, interface{}) interface{}, ctx interface{}) BoundFunction {
	return BoundFunction{apply, ctx}
}
