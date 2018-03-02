package itrex

import (
	"fmt"
	"strings"
	"container/list"
	"github.com/basking2/sdsai-itrex-golang/pkg/sdsai/iterator"
)

func ToBool(v interface{}) bool {
	switch b := v.(type) {
	case bool:
		return b
	case string:
		s := strings.ToLower(b)
		return s == "true" || s == "t" || s == "1"
	case int8:
		return b != 0
	case int16:
		return b != 0
	case int32:
		return b != 0
	case int64:
		return b != 0
	case uint8:
		return b != 0
	case uint16:
		return b != 0
	case uint32:
		return b != 0
	case uint64:
		return b != 0
	case float32:
		return b != 0
	case float64:
		return b != 0
	default:
		return false
	}
}

func ToString(v interface{}) string {
	switch s := v.(type) {
	case bool:
		if s {
			return "true"
		} else {
			return "false"
		}
	case string:
		return s
	case int8:
		return fmt.Sprintf("%d", s)
	case int16:
		return fmt.Sprintf("%d", s)
	case int32:
		return fmt.Sprintf("%d", s)
	case int64:
		return fmt.Sprintf("%d", s)
	case uint8:
		return fmt.Sprintf("%d", s)
	case uint16:
		return fmt.Sprintf("%d", s)
	case uint32:
		return fmt.Sprintf("%d", s)
	case uint64:
		return fmt.Sprintf("%d", s)
	case float32:
		return fmt.Sprintf("%f", s)
	case float64:
		return fmt.Sprintf("%f", s)
	default:
		return fmt.Sprintf("%x", s)
	}
}

func ToIterator(v interface{}) iterator.Iterator {
	switch v2 := v.(type) {
	case iterator.Iterator:
		return v2
	case *list.List:
		return iterator.NewListIterator(v2)
	case []interface{}:
		return iterator.NewArrayIterator(v2)
	default:
		return nil
	}
}
