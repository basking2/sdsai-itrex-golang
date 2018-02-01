package itrex

import (
	"github.com/basking2/sdsai-itrex-golang/sdsai/iterator"
)

type Context struct {
	parent           *Context
	functionRegistry map[string]FunctionInterface
	environment      map[string]interface{}
	arguments        iterator.Iterator
}

func NewContext() *Context {
	return ChildContext(nil)
}

func ChildContext(parent *Context) *Context {
	c := Context{
		parent,
		make(map[string]FunctionInterface),
		make(map[string]interface{}),
		iterator.EmptyIterator{},
	}

	return &c
}

func (parent *Context) FunctionCall(arguments iterator.Iterator) *Context {
	c := Context{
		parent,
		parent.functionRegistry,
		parent.environment,
		arguments,
	}

	return &c
}

func (c *Context) Set(name string, val interface{}) {
	c.environment[name] = val
}

func (c *Context) Update(name string, val interface{}) {
	for ths := c; ths != nil; ths = ths.parent {
		if _, ok := ths.environment[name]; ok {
			ths.environment[name] = val
		}
	}

	panic("Not found.")
}

func (c *Context) Get(name string) interface{} {
	for ths := c; ths != nil; ths = ths.parent {
		if val, ok := ths.environment[name]; ok {
			return val
		}
	}

	return nil
}

func (c *Context) ContainsKey(name string) bool {
	for ths := c; ths != nil; ths = ths.parent {
		if _, ok := ths.environment[name]; ok {
			return true
		}
	}

	return false
}

func (c *Context) Register(name string, f FunctionInterface) {
	c.functionRegistry[name] = f
}

func (c *Context) GetFunction(name string) FunctionInterface {
	for ths := c; ths != nil; ths = ths.parent {
		if val, ok := ths.functionRegistry[name]; ok {
			return val
		}
	}

	return nil
}

func (c *Context) SetArguments(args iterator.Iterator) {
	c.arguments = args
}

func (c *Context) GetArguments() iterator.Iterator {
	return c.arguments
}
