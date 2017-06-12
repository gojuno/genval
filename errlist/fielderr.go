package errlist

import (
	"fmt"

	"github.com/pkg/errors"
)

const UnknownField = "unknown"

type Field struct {
	Field string `json:"field"`
	Err   error  `json:"err"`
}

// Error implements `error` interface
func (e Field) Error() string {
	return fmt.Sprintf("%s: %v", e.Field, e.Err)
}

func (e Field) MarshalJSON() ([]byte, error) {
	data := fmt.Sprintf(`{"field":%q,"error":%q}`, e.Field, e.Err.Error())
	return []byte(data), nil
}

func NewField(field, message string, args ...interface{}) error {
	return &Field{Field: field, Err: errors.Errorf(message, args...)}
}
