package iterator

type ConcatinatedIterator struct {
	iterators []Iterator
	i         int
}

// Position ci.i at the current or next iterator for which HasNext() returns true for.
func (ci *ConcatinatedIterator) nextI() {
	// Position i at the first iterator that has a next element.
	for ci.i < len(ci.iterators) && !ci.iterators[ci.i].HasNext() {
		ci.i += 1
	}
}

func NewConcatinatedIterator(iterators ...Iterator) *ConcatinatedIterator {
	ci := ConcatinatedIterator{iterators, 0}

	ci.nextI()

	return &ci
}

func (i *ConcatinatedIterator) HasNext() bool {
	return i.i < len(i.iterators)
}

func (ci *ConcatinatedIterator) Next() interface{} {
	v := ci.iterators[ci.i].Next()

	ci.nextI()

	return v
}
