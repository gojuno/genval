package types

import (
	"fmt"
	"go/ast"
	"io"
)

func NewInterface() *typeInterface {
	return &typeInterface{}
}

type typeInterface struct {
	funcs []string
}

func (t typeInterface) Type() string {
	return "interface"
}

func (t *typeInterface) SetValidateTag(tag ValidatableTag) error {
	switch tag.Key() {
	case InterfaceFuncKey:
		for _, v := range parseFuncsParam(tag.(SimpleTag).Param) {
			t.funcs = append(t.funcs, v)
		}
	default:
		return ErrUnusedTag
	}
	return nil
}

func (t typeInterface) Generate(w io.Writer, cfg GenConfig, name Name) {
	for _, f := range t.funcs {
		fmt.Fprintf(w, "if err := %s(%s); err != nil {\n", f, name.Full())
		fmt.Fprintf(w, "	errs.AddField(%s, err)\n", name.LabelName())
		fmt.Fprintf(w, "}\n")
	}
}

func (t typeInterface) Validate() error {
	return nil
}

func (t typeInterface) Expr() ast.Expr {
	return nil
}
