package itrml

import (
  "github.com/basking2/sdsai-itrex-golang/pkg/sdsai/iterator"
  "container/list"
  "errors"
  "fmt"
)

type StringWriter interface {
	WriteString(s string) (int, error)
}

func ToString(expr interface{}, w StringWriter) (int, error) {
	total := 0

	switch expr := expr.(type) {
	case string:
		n, err := w.WriteString("\"" + expr + "\"")
		if err != nil {
			return 0, err
		}
		total += n
	case int:
		s := fmt.Sprintf("%d", expr)
		n, err := w.WriteString(s)
		if err != nil {
			return 0, err
		}
		total += n
	case int32:
		s := fmt.Sprintf("%d", expr)
		n, err := w.WriteString(s)
		if err != nil {
			return 0, err
		}
		total += n
	case int64:
		s := fmt.Sprintf("%d", expr)
		n, err := w.WriteString(s + "L")
		if err != nil {
			return 0, err
		}
		total += n
	case uint:
		s := fmt.Sprintf("%d", expr)
		n, err := w.WriteString(s)
		if err != nil {
			return 0, err
		}
		total += n
	case uint32:
		s := fmt.Sprintf("%d", expr)
		n, err := w.WriteString(s)
		if err != nil {
			return 0, err
		}
		total += n
	case uint64:
		s := fmt.Sprintf("%d", expr)
		n, err := w.WriteString(s + "L")
		if err != nil {
			return 0, err
		}
		total += n
	case float32:
		s := fmt.Sprintf("%f", expr)
		n, err := w.WriteString(s + "F")
		if err != nil {
			return 0, err
		}
		total += n
	case float64:
		s := fmt.Sprintf("%f", expr)
		n, err := w.WriteString(s + "D")
		if err != nil {
			return 0, err
		}
		total += n
	case bool:
		if expr {
			n, err := w.WriteString("[boolean true]")
			if err != nil {
				return 0, err
			}
			total += n
		} else {
			n, err := w.WriteString("[boolean false]")
			if err != nil {
				return 0, err
			}
			total += n
		}
	case iterator.Iterator:
		n, err := w.WriteString("[ ")
		if err != nil {
			return 0, err
		}
		total += n
		for expr.HasNext() {
			e := expr.Next()
			n, err = ToString(e, w)
			if err != nil {
				return 0, err
			}
			total += n
			n, err = w.WriteString(" ")
			if err != nil {
				return 0, err
			}
			total += n
		}
		n, err = w.WriteString(" ]")
		if err != nil {
			return 0, err
		}
		total += n
	case *list.List:
		n, err := w.WriteString("[ ")
		if err != nil {
			return 0, err
		}
		total += n
		for e := expr.Front(); e != nil; e = e.Next() {
			n, err = ToString(e.Value, w)
			if err != nil {
				return 0, err
			}
			total += n
			n, err = w.WriteString(" ")
			if err != nil {
				return 0, err
			}
			total += n
		}
		n, err = w.WriteString(" ]")
		if err != nil {
			return 0, err
		}
		total += n
	default:
		return 0, errors.New("Unhandeld type.")
	}

	// Happy Exit.
	return total, nil
}
