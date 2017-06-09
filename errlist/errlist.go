package errlist

import (
	"bytes"
	"fmt"
)

type ErrList []error

// Error implements `error` interface
func (e *ErrList) Error() string {
	var buffer bytes.Buffer

	buffer.WriteString("[")
	for i, err := range *e {
		buffer.WriteString(err.Error())
		if i < len(*e)-1 {
			buffer.WriteString(", ")
		}
	}
	buffer.WriteString("]")

	return buffer.String()
}

func (e *ErrList) Add(err error) *ErrList {
	if err == nil {
		return e
	}

	if errs, ok := err.(*ErrList); ok {
		*e = append(*e, *errs...)
		return e
	}

	*e = append(*e, err)
	return e
}

func (e *ErrList) Addf(msg string, args ...interface{}) *ErrList {
	return e.Add(fmt.Errorf(msg, args...))
}

func (e *ErrList) AddFieldErrf(field, msg string, args ...interface{}) *ErrList {
	return e.Add(NewFieldErr(field, msg, args...))
}

func (e *ErrList) AddFieldErr(field string, err error) *ErrList {
	return e.Add(&FieldErr{Field: field, Err: err})
}
