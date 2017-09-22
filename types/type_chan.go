package types

import (
	"go/ast"
	"io"
)

func NewChan() *typeChan {
	return &typeChan{}
}

type typeChan struct {
}

func (t typeChan) Type() string {
	return "chan"
}

func (t *typeChan) SetValidateTag(tag ValidatableTag) error {
	return ErrUnusedTag
}

func (t typeChan) Generate(w io.Writer, cfg GenConfig, name Name) {

}

func (t typeChan) Validate() error {
	return nil
}

func (t typeChan) Expr() ast.Expr {
	return nil
}
