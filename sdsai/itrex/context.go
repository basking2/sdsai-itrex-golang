package itrex

type Context struct {
	parent *Context
	functionRegistry map[string]func(*Iterator, *Context)interface{}
	environment map[string]interface{}
	arguments Iterator
}

func NewContext() *Context {
	return ChildContext(nil)
}

func ChildContext(parent *Context) *Context {
	c := Context{
		parent,
		make(map[string]func(*Iterator, *Context)interface{}),
		make(map[string]interface{}),
		EmptyIterator{},
	}

	return &c
}

func (parent *Context) FunctionCall(arguments Iterator) *Context {
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

func (c *Context) Register(name string, f func(*Iterator, *Context)interface{}) {
	c.functionRegistry[name] = f
}

func (c *Context) GetFunction(name string) func(*Iterator, *Context)interface{} {
	for ths := c; ths != nil; ths = ths.parent {
		if val, ok := ths.functionRegistry[name]; ok {
			return val
		}
	}

	return nil
}

func (c *Context) SetArguments(args Iterator) {
	c.arguments = args
}

func (c *Context) GetArguments() Iterator {
	return c.arguments
}
