package errlist

import (
	"bytes"
	"fmt"
)

type ErrList []error

func (e *ErrList) Error() string {
	var buffer bytes.Buffer

	buffer.WriteString("[")
	for _, err := range *e {
		buffer.WriteString(fmt.Sprintf(`'%v' `, err))
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
