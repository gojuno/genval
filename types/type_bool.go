package types

import "io"

const Bool string = "bool"

func NewBool() *typeBool {
	return &typeBool{}
}

type typeBool struct {
}

func (t typeBool) Type() string {
	return Bool
}

func (t *typeBool) SetValidateTag(tag ValidatableTag) error {
	return ErrUnusedTag
}

func (t typeBool) Generate(w io.Writer, cfg GenConfig, name Name) {

}

func (t typeBool) Validate() error {
	return nil
}
