package itrex

import (
	"container/list"
	"errors"
	"fmt"
	"github.com/basking2/sdsai-itrex-golang/pkg/sdsai/iterator"
)

type PrintFunction struct {
	Out interface {
		WriteString(string) (int, error)
	}
}

func (f *PrintFunction) Apply(i iterator.Iterator, c *Context) interface{} {

	l := list.New()

	for i.HasNext() {
		v := i.Next()
		l.PushBack(v)
		switch v := v.(type) {
		case string:
			f.Out.WriteString(string(v) + ":string\n")
		case int32:
			f.Out.WriteString(string(v) + ":int32\n")
		case int64:
			f.Out.WriteString(string(v) + ":int64\n")
		case uint64:
			f.Out.WriteString(string(v) + ":uint64\n")
		case float32:
			f.Out.WriteString(fmt.Sprintf("%f:float32\n", v))
		case float64:
			f.Out.WriteString(fmt.Sprintf("%f:float64\n", v))
		case bool:
			if v {
				f.Out.WriteString("true:bool\n")
			} else {
				f.Out.WriteString("false:bool\n")
			}
		default:
			f.Out.WriteString(fmt.Sprintf("%x:ptr\n", v))
		}
	}

	return l
}

type TraceFunction struct {
	Out interface {
		WriteString(string) (int, error)
	}
}

func (f *TraceFunction) Apply(i iterator.Iterator, c *Context) interface{} {

	l := list.New()

	var itrexFunction FunctionInterface

	// First get the function.
	switch fptr := i.Next().(type) {
	case FunctionInterface:
		f.Out.WriteString(fmt.Sprintf("[ %x", fptr))
		itrexFunction = fptr
	case string:
		f.Out.WriteString(fmt.Sprintf("[ %s", fptr))
		itrexFunction = c.GetFunction(fptr)
	default:
		f.Out.WriteString(fmt.Sprintf("Could not find function: %x.", fptr))
		return nil
	}

	for i.HasNext() {
		v := i.Next()
		l.PushBack(v)
		f.Out.WriteString(" " + ToString(v))
	}

	f.Out.WriteString(" ]")

	// Redefine i and c for a function call.
	i = iterator.NewListIterator(l)
	return itrexFunction.Apply(i, c)
}

type IfFunction struct{}

func (IfFunction) Apply(i iterator.Iterator, c *Context) interface{} {

	if !i.HasNext() {
		return errors.New("No predicate.")
	}

	var boolVal = ToBool(i.Next())

	if boolVal {
		if i.HasNext() {
			return i.Next()
		} else {
			return errors.New("No then clause in if.")
		}
	} else {
		if !i.HasNext() {
			return errors.New("No then or elses clause in if.")
		}

		switch v := i.(type) {
		case *EvaluatingIterator:
			v.EvaluationEnabled = false
			v.Next()
			v.EvaluationEnabled = true
		default:
			v.Next()
		}

		if !i.HasNext() {
			return errors.New("No else clause in if.")
		}

		return i.Next()
	}
}

type SetFunction struct{}

func (SetFunction) Apply(i iterator.Iterator, c *Context) interface{} {
	if i.HasNext() {
		name := ToString(i.Next())
		if i.HasNext() {
			v := i.Next()
			c.Set(name, v)
			return v
		} else {
			c.Set(name, nil)
			return nil
		}
	} else {
		return nil
	}
}

type UpdateFunction struct{}

func (UpdateFunction) Apply(i iterator.Iterator, c *Context) interface{} {
	if i.HasNext() {
		name := ToString(i.Next())
		if i.HasNext() {
			v := i.Next()
			c.Update(name, v)
			return v
		} else {
			c.Update(name, nil)
			return nil
		}
	} else {
		return nil
	}
}

type GetFunction struct{}

func (GetFunction) Apply(i iterator.Iterator, c *Context) interface{} {
	if i.HasNext() {
		name := ToString(i.Next())
		return c.Get(name)
	} else {
		return nil
	}
}

type LastFunction struct{}

func (LastFunction) Apply(i iterator.Iterator, c *Context) interface{} {
	var v interface{} = nil
	for i.HasNext() {
		v = i.Next()
	}
	return v
}

type LetFunction struct{}

func (f LetFunction) Apply(i iterator.Iterator, c *Context) interface{} {

	switch evaluatingContext := i.(type) {
	case *EvaluatingIterator:
		newContext := c.LetContext()
		var v interface{} = nil
		for evaluatingContext.HasNext() {
			v = evaluatingContext.NextCtx(newContext)
		}
		return v
	default:
		return errors.New("Can only use let with an EvaluatingIterator as input.")
	}
}

type MapFunction struct{}

func (f MapFunction) Apply(i iterator.Iterator, c *Context) interface{} {
	if ! i.HasNext() {
		return nil
	}

	var fi FunctionInterface = nil

	switch funcInterface := i.Next().(type) {
	case FunctionInterface:
		fi = funcInterface
		if fi == nil {
			return errors.New("Function to map was not found.")
		}
	case string:
		fi = c.GetFunction(funcInterface)
		if fi == nil {
			return errors.New("Function to map was not found: " + funcInterface)
		}
	default:
		return errors.New("First argument is not a function.")
	}

	if i.HasNext() {
		v := ToIterator(i.Next())
		if v == nil {
			return errors.New("Second argument is not iterable.")
		} else {
			return iterator.NewMappingIterator(
				v,
				func(i interface{}) interface{} {
					v := make([]interface{}, 1)
					v[0] = i
					return fi.Apply(iterator.NewArrayIterator(v), c)
				})
		}
	} else {
		return errors.New("There is no second iterable argument.")
	}

	return errors.New("Map Function: Second argument. No an iterator.")
}

type FunctionFunction struct{
	evaluator *Evaluator
}

func (f FunctionFunction) Apply(i iterator.Iterator, c *Context) interface{} {
	var functionBody interface{}

	switch itr := i.(type) {
	case *EvaluatingIterator:
		itr.EvaluationEnabled = false
		functionBody = itr.Next()
		itr.EvaluationEnabled = true
	default:
		functionBody = itr.Next()
	}

	return NewBoundFunction(
		func(args iterator.Iterator, c *Context, functionBody interface{}) interface{} {
			return f.evaluator.Evaluate(functionBody, c.FunctionCall(args))
		},
		functionBody)
}

type FnFunction struct {
	evaluator *Evaluator
}

func (f FnFunction) Apply(i iterator.Iterator, c *Context) interface{} {

	if functionName := ToString(i.Next()); functionName != "" {
		if ! i.HasNext() {
			return c.GetFunction(functionName)
		}

		var functionBody interface{}

		switch itr := i.(type) {
		case *EvaluatingIterator:
			itr.EvaluationEnabled = false
			functionBody = itr.Next()
			itr.EvaluationEnabled = true
		default:
			functionBody = itr.Next()
		}

		boundFunction := NewBoundFunction(
			func(args iterator.Iterator, c *Context, functionBody interface{}) interface{} {
				return f.evaluator.Evaluate(functionBody, c.FunctionCall(args))
			},
			functionBody)

		c.Register(functionName, boundFunction)

		return boundFunction
	} else {
		return nil
	}
}
