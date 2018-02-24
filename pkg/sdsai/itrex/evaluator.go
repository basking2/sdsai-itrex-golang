package itrex

import (
	"os"
	"container/list"
	"github.com/basking2/sdsai-itrex-golang/pkg/sdsai/iterator"
)

type Evaluator struct {
	RootContext *Context
}

func NewEvaluator() *Evaluator {
	e := Evaluator{NewContext()}

	e.Register("print", &PrintFunction{os.Stdout})
	e.Register("printErr", &PrintFunction{os.Stderr})

	return &e
}

func (e *Evaluator) Register(name string, f FunctionInterface) {
	e.RootContext.Register(name, f)
}

func (e *Evaluator) Evaluate(o interface{}, context *Context) interface{} {
	switch o2 := o.(type) {
	case *EvaluatingIterator:
		return e.EvaluateEvaluatingIterator(o2)
	case iterator.Iterator:
		ei := EvaluatingIterator{o2, true, e, context}
		return e.EvaluateEvaluatingIterator(&ei)
	case *list.List:
		ai := iterator.NewListIterator(o2)
		ei := EvaluatingIterator{ai, true, e, context}
		return e.EvaluateEvaluatingIterator(&ei)
	case []interface{}:
		ai := iterator.NewArrayIterator(o2)
		ei := EvaluatingIterator{ai, true, e, context}
		return e.EvaluateEvaluatingIterator(&ei)
	default:
		return o
	}
}

func (e *Evaluator) EvaluateEvaluatingIterator(ei *EvaluatingIterator) interface{} {
	if ei.HasNext() {
		return &iterator.EmptyIterator{}
	}

	operatorObject := ei.Next()
	var operator FunctionInterface

	switch operatorObjectT := operatorObject.(type) {
	case FunctionInterface:
		operator = operatorObjectT
	case string:
		operator = ei.context.GetFunction(operatorObjectT)
	default:
		panic("Cannot handle type.")
	}

	return operator.Apply(ei, ei.context)
}
