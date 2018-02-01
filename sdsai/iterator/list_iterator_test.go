package iterator

import (
	"testing"
  "container/list"
)

func TestSimpleListWalk(t *testing.T) {
	var data = list.New()
  data.PushBack("println")
  data.PushBack(3)
	ai := NewListIterator(data)
	ai.HasNext()
	if ai.Next().(string) != "println" {
		t.Error("Expected value of println.")
	}

	ai.HasNext()
	if v := ai.Next().(int); v != 3 {
		t.Error("Expected value of 3.")
	}

	if ai.HasNext() {
		t.Error("Iterator did not terminate as expected.")
	}
}
