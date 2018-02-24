package itrex

import (
  "testing"
  "github.com/basking2/sdsai-itrex-golang/pkg/sdsai/itrml"
)

func TestPrintFunction(t *testing.T) {
  e := NewEvaluator()

  println("----------------------")
  expr, err := itrml.ParseExpression("[print hi how are you]")
  if err != nil {
    panic(err.Error())
  }
  e.Evaluate(expr, e.RootContext)
  println("----------------------")

  panic("This is not correct.")
}
