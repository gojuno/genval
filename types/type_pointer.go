package types

import (
	"fmt"
	"io"
)

func NewPointer(inner TypeDef) *typePointer {
	return &typePointer{
		nullable:  true,
		innerType: inner,
	}
}

type typePointer struct {
	nullable  bool
	innerType TypeDef
}

func (t typePointer) Type() string {
	return t.innerType.Type()
}

func (t *typePointer) SetTag(tag Tag) error {
	switch tag.Key() {
	case PointerNullableKey:
		t.nullable = true
	case PointerNotNullKey:
		t.nullable = false
	default:
		return t.innerType.SetTag(tag)
	}
	return nil
}

func (t typePointer) Generate(w io.Writer, cfg GenConfig, name Name) {
	if t.nullable {
		fmt.Fprintf(w, "if %s != nil {\n", name.Full())
		t.innerType.Generate(w, cfg, name.WithPointer())
		fmt.Fprintf(w, "}\n")
	} else {
		fmt.Fprintf(w, "if %s == nil {\n", name.Full())
		fmt.Fprintf(w, "    errs.AddFieldErrf(%s, \"cannot be nil\")\n", name.LabelName())
		fmt.Fprintf(w, "}\n")
		t.innerType.Generate(w, cfg, name.WithPointer())
	}
}

func (t typePointer) Validate() error {
	return t.innerType.Validate()
}
