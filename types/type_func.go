package types

import (
	"go/ast"
	"io"
)

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

func (t typeFunc) Expr() ast.Expr {
	return nil
}
