package errlist

import (
	"bytes"

	"github.com/pkg/errors"
)

type List []error

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
		*e = append(*e, err)
		return e
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

	switch v := err.(type) {
	case List:
		for _, err := range v {
			if fieldErr, ok := err.(Field); ok {
				*e = append(*e, Field{Field: field + "." + fieldErr.Field, Err: fieldErr.Err})
			} else {
				*e = append(*e, err)
			}
		}
	case Field:
		*e = append(*e, errors.Errorf("%s.%v", field, v.Error()))
	default:
		*e = append(*e, Field{Field: field, Err: err})
	}

	return e
}

func (e List) ErrorOrNil() error {
	if len(e) == 0 {
		return nil
	}

	if len(e) == 1 {
		return e[0]
	}

	return e
}
