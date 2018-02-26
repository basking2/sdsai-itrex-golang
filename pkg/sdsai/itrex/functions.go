package itrex

import (
  "fmt"
  "errors"
  "container/list"
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
      f.Out.WriteString(string(v)+":string\n")
    case int32:
      f.Out.WriteString(string(v)+":int32\n")
    case int64:
      f.Out.WriteString(string(v)+":int64\n")
    case uint64:
      f.Out.WriteString(string(v)+":uint64\n")
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

type LetFunction struct{
  evaluator *Evaluator
}

func (f LetFunction) Apply(i iterator.Iterator, c *Context) interface{} {
  childContext := ChildContext(c)
  f.evaluator.RootContext = childContext
  defer func() {
    f.evaluator.RootContext = f.evaluator.RootContext.parent
  }()

  var v interface{} = nil
  for i.HasNext() {
    v = i.Next()
  }
  return v
}
