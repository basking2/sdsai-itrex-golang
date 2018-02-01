package itrex

type EmptyIterator struct{}

func (e EmptyIterator) HasNext() bool {
	return false
}

func (e EmptyIterator) Next() interface{} {
	panic("No next value.")
}
