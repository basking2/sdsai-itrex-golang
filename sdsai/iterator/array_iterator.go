package iterator

type ArrayIterator struct {
	data  []interface{}
	index int
}

func NewArrayIterator(data []interface{}) *ArrayIterator {
	return &ArrayIterator{data, 0}
}

func (a *ArrayIterator) HasNext() bool {
	return a.index < len(a.data)
}

func (a *ArrayIterator) Next() interface{} {
	r := a.data[a.index]
	a.index += 1
	return r
}
