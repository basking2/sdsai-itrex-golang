package itrex

type Evaluator struct {
	RootContext *Context
}

func (e *Evaluator) Register(name string, f FunctionInterface) {
	e.RootContext.Register(name, f)
}

func (e *Evaluator) Evaluate(o interface{}, context *Context) interface{} {
  switch o2 := o.(type) {
  case *EvaluatingIterator:
    return e.EvaluateEvaluatingIterator(o2)
  case Iterator:
    ei := EvaluatingIterator{o2, true, e, context}
    return e.EvaluateEvaluatingIterator(&ei)
  default:
    return o
  }
}

func (e *Evaluator) EvaluateEvaluatingIterator(ei *EvaluatingIterator) interface{} {
  if (ei.HasNext()) {
    return &EmptyIterator{}
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

  return operator(ei, ei.context)
}
