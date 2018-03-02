package iterator

type MappingIterator struct {
  i Iterator
  f func (interface{}) interface{}
}

func NewMappingIterator(i Iterator, f func(interface{}) interface{}) *MappingIterator {
  return &MappingIterator{i, f}
}

func (m *MappingIterator) HasNext() bool {
  return m.i.HasNext()
}

func (m *MappingIterator) Next() interface{} {
  n := m.Next()
  return m.f(n)
}
