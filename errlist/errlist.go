package errlist

import (
	"bytes"

	"github.com/pkg/errors"
)

type List []Field

// Error implements `error` interface
func (e List) Error() string {
	var buffer bytes.Buffer

	buffer.WriteString("[")
	for i, err := range e {
		buffer.WriteString(err.Error())
		if i < len(e)-1 {
			buffer.WriteString(", ")
		}
	}
	buffer.WriteString("]")

	return buffer.String()
}

func (e *List) Add(err error) *List {
	if err == nil {
		return e
	}

	list, ok := err.(List)
	if !ok {
		return e.AddField(UnknownField, err)
	}
	for _, fe := range list {
		*e = append(*e, fe)
	}
	return e
}

func (e *List) AddFieldf(field, msg string, args ...interface{}) *List {
	return e.AddField(field, errors.Errorf(msg, args...))
}

func (e *List) AddField(field string, err error) *List {
	if err == nil {
		return e
	}

	switch errTyped := err.(type) {
	case List:
		for _, childErr := range errTyped {
			*e = append(*e, Field{Field: field + "." + childErr.Field, Err: childErr.Err})
		}
	case error:
		*e = append(*e, Field{Field: field, Err: err})
	}

	return e
}

func (e List) HasErrors() bool {
	return len(e) > 0
}

func (e List) ErrorOrNil() error {
	if e.HasErrors() {
		return e
	}
	return nil
}
