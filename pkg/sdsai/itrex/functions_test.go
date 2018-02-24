package itrex

import (
  "testing"
  "bytes"
  "github.com/basking2/sdsai-itrex-golang/pkg/sdsai/itrml"
)

func TestPrintFunction(t *testing.T) {
  e := NewEvaluator()
  buffer := bytes.Buffer{}
  e.Register("print", &PrintFunction{&buffer})

  expr, err := itrml.ParseExpression("[print hi how are you]")
  if err != nil {
    panic(err.Error())
  }
  e.Evaluate(expr, e.RootContext)
  if buffer.String() != "hi:string\nhow:string\nare:string\nyou:string\n" {
    t.Error("Expected \"hihowareyou\" but got " + buffer.String())
  }

}
