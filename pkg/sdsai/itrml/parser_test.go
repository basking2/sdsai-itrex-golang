package itrml

import (
	"container/list"
	"testing"
)

func TestParseAList(t *testing.T) {

	switch v, _ := ParseExpression("[]"); v.(type) {
	case *list.List:
	default:
		t.Error("Expected list.")
	}
}

func TestParseAList2(t *testing.T) {

	switch v, _ := ParseExpression("    []    "); v.(type) {
	case *list.List:
	default:
		t.Error("Expected list.")
	}
}

func TestParseAListWithData(t *testing.T) {

	switch v, _ := ParseExpression("[string1 \"string two\" 1l 2d 3]"); l := v.(type) {
	case *list.List:
		f := l.Front()
		if f.Value.(string) != "string1" {
			t.Error("Expected string1")
		}

		f = f.Next()
		if f.Value.(string) != "string two" {
			t.Errorf("Expected string two. Got %s", f.Value)
		}

		f = f.Next()
		if f.Value.(int64) != 1 {
			t.Errorf("Expected int of 1")
		}

		f = f.Next()
		if f.Value.(float64) != 2.0 {
			t.Errorf("Expected int of 2.0.")
		}

		f = f.Next()
		if f.Value.(int64) != 3 {
			t.Errorf("Expected int of 3.")
		}
	default:
		t.Error("Expected list.")
	}
}

func TestQuotes(t *testing.T) {
	switch v, _ := ParseExpression("[\"He said, \\\\\\\"Hello.\\\"\"]"); l := v.(type) {
	case *list.List:
		f := l.Front()
		if f.Value.(string) != "He said, \\\"Hello.\"" {
			t.Error("Expected string1")
		}
	default:
		t.Error("Expected list.")
	}
}
