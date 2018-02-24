package itrex

import (
  "fmt"
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
