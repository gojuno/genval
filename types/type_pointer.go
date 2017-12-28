package types

import (
	"fmt"
	"go/ast"
	"io"
)

func NewPointer(inner TypeDef) *TypePointer {
	return &TypePointer{
		nullable:  true,
		innerType: inner,
	}
}

type TypePointer struct {
	nullable  bool
	innerType TypeDef
}

func (t *TypePointer) Type() string {
	return t.innerType.Type()
}

func (t *TypePointer) SetInnerType(newType TypeDef) *TypePointer {
	t.innerType = newType
	return t
}

func (t *TypePointer) SetValidateTag(tag ValidatableTag) error {
	switch tag.Key() {
	case PointerNullableKey:
		t.nullable = true
	case PointerNotNullKey:
		t.nullable = false
	default:
		return t.innerType.SetValidateTag(tag)
	}
	return nil
}

func (t TypePointer) NeedGenerate() bool {
	return true
}

func (t *TypePointer) Generate(w io.Writer, cfg GenConfig, name Name) {
	if t.nullable {
		fmt.Fprintf(w, "if %s != nil {\n", name.Full())
		t.innerType.Generate(w, cfg, name.WithPointer())
		fmt.Fprintf(w, "}\n")
	} else {
		fmt.Fprintf(w, "if %s == nil {\n", name.Full())
		fmt.Fprintf(w, "    errs.AddFieldf(%s, \"cannot be nil\")\n", name.LabelName())
		fmt.Fprintf(w, "} else {\n")
		t.innerType.Generate(w, cfg, name.WithPointer())
		fmt.Fprintf(w, "}\n")
	}
}

func (t *TypePointer) Validate() error {
	return t.innerType.Validate()
}

func (t *TypePointer) Expr() ast.Expr {
	return t.innerType.Expr()
}
