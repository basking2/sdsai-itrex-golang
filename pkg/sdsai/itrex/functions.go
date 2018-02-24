package itrex

import (
  "fmt"
  "github.com/basking2/sdsai-itrex-golang/pkg/sdsai/iterator"
)

type PrintFunction struct {
  Out interface {
    WriteString(string) (int, error)
  }
}

func (f *PrintFunction) Apply(i iterator.Iterator, c *Context) interface{} {
  for i.HasNext() {
    switch v := i.Next().(type) {
    case string:
      f.Out.WriteString(string(v))
    case int32:
      f.Out.WriteString(string(v))
    case int64:
      f.Out.WriteString(string(v))
    case uint64:
      f.Out.WriteString(string(v))
    case float32:
      f.Out.WriteString(fmt.Sprintf("%f", v))
    case float64:
      f.Out.WriteString(fmt.Sprintf("%f", v))
    case bool:
      if v {
        f.Out.WriteString("true")
      } else {
        f.Out.WriteString("false")
      }
    default:
      f.Out.WriteString(fmt.Sprintf("%x", v))
    }
  }

  return nil
}
