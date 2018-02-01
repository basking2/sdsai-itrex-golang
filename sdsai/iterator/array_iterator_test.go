package iterator

import (
	"testing"
)

func TestSimpleWalk(t *testing.T) {
	var data = make([]interface{}, 2)
	data[0] = "println"
	data[1] = 3
	ai := ArrayIterator{data, 0}
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
