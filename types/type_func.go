package types

import "io"

func NewFunc() *typeFunc {
	return &typeFunc{}
}

type typeFunc struct {
}

func (t typeFunc) Type() string {
	return "func"
}

func (t *typeFunc) SetValidateTag(tag ValidatableTag) error {
	return ErrUnusedTag
}

func (t typeFunc) Generate(w io.Writer, cfg GenConfig, name Name) {

}

func (t typeFunc) Validate() error {
	return nil
}
