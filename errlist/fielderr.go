package errlist

import "fmt"

type FieldErr struct {
	Field string `json:"field"`
	Err   error  `json:"err"`
}

// Error implements `error` interface
func (e *FieldErr) Error() string {
	return fmt.Sprintf("%s: %v", e.Field, e.Err)
}

func NewFieldErr(field, message string, args ...interface{}) error {
	return &FieldErr{Field: field, Err: fmt.Errorf(message, args...)}
}
