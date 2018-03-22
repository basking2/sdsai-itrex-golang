package itrex

import (
	"container/list"
	"github.com/basking2/sdsai-itrex-golang/pkg/sdsai/iterator"
	"os"
)

type Evaluator struct {
	RootContext *Context
}

func NewEvaluator() *Evaluator {
	e := Evaluator{NewContext()}

	e.Register("print", &PrintFunction{os.Stdout})
	e.Register("printErr", &PrintFunction{os.Stderr})
	e.Register("trace", &TraceFunction{os.Stdout})
	e.Register("traceErr", &TraceFunction{os.Stderr})
	e.Register("if", IfFunction{})
	e.Register("set", SetFunction{})
	e.Register("update", UpdateFunction{})
	e.Register("get", GetFunction{})
	e.Register("let", LetFunction{})
	e.Register("last", LastFunction{})
	e.Register("function", FunctionFunction{&e})
	e.Register("register", RegisterFunction{})
	e.Register("args", ArgsFunction{})
	e.Register("arg", ArgFunction{})
	e.Register("hasArg", HasArgFunction{})
	e.Register("fn", FnFunction{&e})
	e.Register("curry", CurryFunction{})
	e.Register("evalItrml", EvalItrMlFunction{&e})
	e.Register("nop", NewBoundFunction(func(i iterator.Iterator, c *Context, cbdata interface{}) interface{} {
		return nil
	}, nil))

	e.Register("boolean", NewBoundFunction(func(i iterator.Iterator, c *Context, cbdata interface{}) interface{} {
		return ToBool(i.Next())
	}, nil))
	e.Register("string", NewBoundFunction(func(i iterator.Iterator, c *Context, cbdata interface{}) interface{} {
		return ToString(i.Next())
	}, nil))
	e.Register("int", NewBoundFunction(func(i iterator.Iterator, c *Context, cbdata interface{}) interface{} {
		return ToInt(i.Next())
	}, nil))
	e.Register("long", NewBoundFunction(func(i iterator.Iterator, c *Context, cbdata interface{}) interface{} {
		return ToLong(i.Next())
	}, nil))
	e.Register("float", NewBoundFunction(func(i iterator.Iterator, c *Context, cbdata interface{}) interface{} {
		return ToFloat(i.Next())
	}, nil))
	e.Register("double", NewBoundFunction(func(i iterator.Iterator, c *Context, cbdata interface{}) interface{} {
		return ToDouble(i.Next())
	}, nil))

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
	if !ei.HasNext() {
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
