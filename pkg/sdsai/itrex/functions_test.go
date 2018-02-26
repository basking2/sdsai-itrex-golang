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

func TestTraceFunction(t *testing.T) {
  e := NewEvaluator()
  buffer := bytes.Buffer{}
  e.Register("trace", &TraceFunction{&buffer})

  expr, err := itrml.ParseExpression("[trace print hi how are you]")
  if err != nil {
    panic(err.Error())
  }
  e.Evaluate(expr, e.RootContext)
  if buffer.String() != "[ print hi how are you ]" {
    t.Error("Expected \"hihowareyou\" but got " + buffer.String())
  }
}

func TestIfFunction(t *testing.T) {
  e := NewEvaluator()

  expr, err := itrml.ParseExpression("[if true 3 4]")
  if err != nil {
    panic(err.Error())
  }
  r := e.Evaluate(expr, e.RootContext)
  if r.(int64) != 3 {
    t.Errorf("%d != 3", r.(int64))
  }
  
  expr, err = itrml.ParseExpression("[if false 3 4]")
  if err != nil {
    panic(err.Error())
  }
  r = e.Evaluate(expr, e.RootContext)
  if r.(int64) != 4 {
    t.Errorf("%d != 4", r.(int64))
  }
}
