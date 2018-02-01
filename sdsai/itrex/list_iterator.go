package itrex

import (
	"container/list"
)

type ListIterator struct {
	// The list.
	data  *list.List

	// The current element in the list.
	// Data from this element has already been returned to the user.
	// In a new ListIterator this will be nil and started will be false.
	// In a finished ListIterator this will be nil and started will be true.
	index *list.Element

	// Has next been called at least once?
	started bool
}

func NewListIterator(l *list.List) *ListIterator {
	li := ListIterator{l, nil, false}
	return &li
}

func (a *ListIterator) HasNext() bool {
	if a.started {
		return a.index.Next() != nil
	} else {
		return a.data.Front() != nil
	}
}

func (a *ListIterator) Next() interface{} {
	if a.started {
		a.index = a.index.Next()
	} else {
		a.started = true
		a.index = a.data.Front()
	}

	return a.index.Value
}
