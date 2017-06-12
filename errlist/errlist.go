package errlist

import (
	"bytes"

	"github.com/pkg/errors"
)

type ErrList []FieldErr

// Error implements `error` interface
func (e ErrList) Error() string {
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

func (e *ErrList) Add(err error) *ErrList {
	if err == nil {
		return e
	}

	list, ok := err.(ErrList)
	if !ok {
		return e.AddFieldErr(UnknownField, err)
	}
	for _, fe := range list {
		*e = append(*e, fe)
	}
	return e
}

func (e *ErrList) AddFieldErrf(field, msg string, args ...interface{}) *ErrList {
	return e.AddFieldErr(field, errors.Errorf(msg, args...))
}

func (e *ErrList) AddFieldErr(field string, err error) *ErrList {
	if err == nil {
		return e
	}

	*e = append(*e, FieldErr{Field: field, Err: err})

	return e
}

func (e ErrList) HasErrors() bool {
	return len(e) > 0
}

func (e ErrList) ErrorOrNil() error {
	if e.HasErrors() {
		return e
	}
	return nil
}
