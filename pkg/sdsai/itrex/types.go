package itrex

import (
	"container/list"
	"fmt"
	"github.com/basking2/sdsai-itrex-golang/pkg/sdsai/iterator"
	"strings"
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

func ToInt(v interface{}) int32 {
	switch b := v.(type) {
	case bool:
		if b {
			return 1
		} else {
			return 0
		}
	case string:
		var i int32
		fmt.Sscanf(b, "%d", &i)
		return i
	case int8:
		return int32(b)
	case int16:
		return int32(b)
	case int32:
		return int32(b)
	case int64:
		return int32(b)
	case uint8:
		return int32(b)
	case uint16:
		return int32(b)
	case uint32:
		return int32(b)
	case uint64:
		return int32(b)
	case float32:
		return int32(b)
	case float64:
		return int32(b)
	default:
		return 0
	}
}

func ToLong(v interface{}) int64 {
	switch b := v.(type) {
	case bool:
		if b {
			return 1
		} else {
			return 0
		}
	case string:
		var i int64
		fmt.Sscanf(b, "%d", &i)
		return i
	case int8:
		return int64(b)
	case int16:
		return int64(b)
	case int32:
		return int64(b)
	case int64:
		return int64(b)
	case uint8:
		return int64(b)
	case uint16:
		return int64(b)
	case uint32:
		return int64(b)
	case uint64:
		return int64(b)
	case float32:
		return int64(b)
	case float64:
		return int64(b)
	default:
		return 0
	}

}

func ToFloat(v interface{}) float32 {
	switch b := v.(type) {
	case bool:
		if b {
			return 1
		} else {
			return 0
		}
	case string:
		var i float32
		fmt.Sscanf(b, "%f", &i)
		return i
	case int8:
		return float32(b)
	case int16:
		return float32(b)
	case int32:
		return float32(b)
	case int64:
		return float32(b)
	case uint8:
		return float32(b)
	case uint16:
		return float32(b)
	case uint32:
		return float32(b)
	case uint64:
		return float32(b)
	case float32:
		return float32(b)
	case float64:
		return float32(b)
	default:
		return 0
	}
}

func ToDouble(v interface{}) float64 {
	switch b := v.(type) {
	case bool:
		if b {
			return 1
		} else {
			return 0
		}
	case string:
		var i float64
		fmt.Sscanf(b, "%f", &i)
		return i
	case int8:
		return float64(b)
	case int16:
		return float64(b)
	case int32:
		return float64(b)
	case int64:
		return float64(b)
	case uint8:
		return float64(b)
	case uint16:
		return float64(b)
	case uint32:
		return float64(b)
	case uint64:
		return float64(b)
	case float32:
		return float64(b)
	case float64:
		return float64(b)
	default:
		return 0
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
