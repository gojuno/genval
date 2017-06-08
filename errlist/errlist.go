package errlist

import (
	"bytes"
	"fmt"
)

type ErrList []error

func (e *ErrList) Error() string {
	var buffer bytes.Buffer

	buffer.WriteString("[")
	for i, err := range *e {
		buffer.WriteString(fmt.Sprintf("'%v'", err))
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
