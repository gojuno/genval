package types

import (
	"go/ast"
	"io"
)

const Byte string = "byte"

func NewByte() *typeByte {
	return &typeByte{}
}

type typeByte struct {
}

func (t typeByte) Type() string {
	return Byte
}

func (t *typeByte) SetValidateTag(tag ValidatableTag) error {
	return ErrUnusedTag
}

func (t typeByte) NeedGenerate() bool {
	return false
}

func (t typeByte) Generate(w io.Writer, cfg GenConfig, name Name) {

}

func (t typeByte) Validate() error {
	return nil
}

func (t typeByte) Expr() ast.Expr {
	return nil
}
