package errlist

import (
	"fmt"

	"github.com/pkg/errors"
)

const UnknownField = "unknown"

type FieldErr struct {
	Field string `json:"field"`
	Err   error  `json:"err"`
}

// Error implements `error` interface
func (e FieldErr) Error() string {
	return fmt.Sprintf("%s: %v", e.Field, e.Err)
}

func (e FieldErr) MarshalJSON() ([]byte, error) {
	data := fmt.Sprintf(`{"field":%q,"error":%q}`, e.Field, e.Err.Error())
	return []byte(data), nil
}

func NewFieldErr(field, message string, args ...interface{}) error {
	return &FieldErr{Field: field, Err: errors.Errorf(message, args...)}
}
