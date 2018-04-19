package itrex

import (
	"container/list"
	"errors"
	"fmt"
	"io/ioutil"
	"github.com/basking2/sdsai-itrex-golang/pkg/sdsai/iterator"
	"github.com/basking2/sdsai-itrex-golang/pkg/sdsai/itrml"
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
	if !i.HasNext() {
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
}

type FunctionFunction struct {
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
		if !i.HasNext() {
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

type CurryFunction struct{}

func (f CurryFunction) Apply(i iterator.Iterator, c *Context) interface{} {
	var functionBody FunctionInterface

	switch arg1 := i.Next().(type) {
	case FunctionInterface:
		functionBody = arg1
	default:
		functionName := ToString(arg1)
		if functionName == "" {
			return nil
		}
		functionBody = c.GetFunction(functionName)
	}

	boundArgs := list.New()
	for i.HasNext() {
		boundArgs.PushBack(i.Next())
	}

	return NewBoundFunction(
		func(args iterator.Iterator, c *Context, cbdata interface{}) interface{} {
			allArgs := iterator.NewConcatinatedIterator(
				iterator.NewListIterator(boundArgs),
				args)
			return functionBody.Apply(allArgs, c)
		},
		nil)
}

type RegisterFunction struct{}
func (f RegisterFunction) Apply(i iterator.Iterator, c *Context) interface{} {
	switch v := i.Next().(type) {
	case string:
		f := i.Next().(FunctionInterface)
		c.Register(v, f)
		return f
	default:
		f := i.Next().(FunctionInterface)
		c.Register(ToString(v), f)
		return f
	}
}

type ArgFunction struct{}
func (f ArgFunction) Apply(i iterator.Iterator, c *Context) interface{} {
	return c.arguments.Next()
}

type ArgsFunction struct{}
func (f ArgsFunction) Apply(i iterator.Iterator, c *Context) interface{} {
	return c.arguments
}

type HasArgFunction struct{}
func (f HasArgFunction) Apply(i iterator.Iterator, c *Context) interface{} {
	return c.arguments.HasNext()
}

type EvalItrMlFunction struct{
	evaluator *Evaluator
}

func (f EvalItrMlFunction) Apply(i iterator.Iterator, c *Context) interface{} {
	file := i.Next().(string)

	exprBytes, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	expr, err := itrml.ParseExpression(string(exprBytes))
	if err != nil {
		return err
	}

	return f.evaluator.Evaluate(expr, c)

}

type NameArgsFunction struct{}
func (f NameArgsFunction) Apply(i iterator.Iterator, c *Context) interface{} {

	var argVal interface{} = nil

	for i.HasNext() && c.arguments.HasNext() {
		v := i.Next().(string)
		argVal = c.arguments.Next()
		c.Set(v, argVal)
	}

	return argVal
}

type CaseFunction struct{}
func (f CaseFunction) Apply(i iterator.Iterator, c *Context) interface{} {
	if (ToBool(i.Next())) {
		v := make([]interface{}, 2)
		v[0] = true
		v[1] = i.Next()
		return v
	} else {
		v := make([]interface{}, 2)
		v[0] = false
		v[1] = nil
		return v
	}
}

type DefaultCaseFunction struct{}
func (f DefaultCaseFunction) Apply(i iterator.Iterator, c *Context) interface{} {
	if (i.HasNext()) {
		v := make([]interface{}, 2)
		v[0] = true
		v[1] = i.Next()
		return v
	} else {
		v := make([]interface{}, 2)
		v[0] = true
		v[1] = true
		return v
	}
}

type CaseListFunction struct{}
func (f CaseListFunction) Apply(i iterator.Iterator, c *Context) interface{} {
	for (i.HasNext()) {
		o := i.Next()
		if (o != nil) {
			itr := ToIterator(o)
			if (itr == nil) {
				if (ToBool(o)) {
					return o
				}
			} else {
				if (itr.HasNext()) {
					if (ToBool(itr.Next())) {
						if (itr.HasNext()) {
							return itr.Next()
						} else {
							return true
						}
					}
				}
			}
		}
	}
	return nil
}

// dict.mk
// dict.put
// dict.get
// callFlattened
// compose
// curry
// foldLeft
// pipeline
// map
// mapFlat
// head
// tail
// list
// listFlatten
// flatten
// flatten2
// string.join
// string.split
// string.concat
// for
// range
// t
// f
// for
// and
// not
// or
// eq
// lt
// lte
// gte
// gt
// version
