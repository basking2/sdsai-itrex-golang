package itrex

import (
	"github.com/basking2/sdsai-itrex-golang/pkg/sdsai/iterator"
)

type EvaluatingIterator struct {
	iterator          iterator.Iterator
	EvaluationEnabled bool
	evaluator         *Evaluator
	context           *Context
}

func (ei *EvaluatingIterator) HasNext() bool {
	return ei.iterator.HasNext()
}

func (ei *EvaluatingIterator) Next() interface{} {
	return ei.NextCtx(ei.context)
}

func (ei *EvaluatingIterator) NextCtx(context *Context) interface{} {
	if ei.EvaluationEnabled {
		return ei.evaluator.Evaluate(ei.iterator.Next(), context)
	} else {
		return ei.iterator.Next()
	}
}
