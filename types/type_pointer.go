package types

import (
	"fmt"
	"io"
)

func NewPointer(inner TypeDef) *typePointer {
	return &typePointer{
		Nullable:  true,
		InnerType: inner,
	}
}

type typePointer struct {
	Nullable  bool
	InnerType TypeDef
}

func (t typePointer) Type() string {
	return t.InnerType.Type()
}

func (t *typePointer) SetTag(tag Tag) error {
	switch tag.Key() {
	case PointerNullableKey:
		t.Nullable = true
	case PointerNotNullKey:
		t.Nullable = false
	default:
		return t.InnerType.SetTag(tag)
	}
	return nil
}

func (t typePointer) Generate(w io.Writer, cfg GenConfig, name Name) {
	if t.Nullable {
		fmt.Fprintf(w, "if %s != nil {\n", name.Full())
		t.InnerType.Generate(w, cfg, name.WithPointer())
		fmt.Fprintf(w, "}\n")
	} else {
		cfg.AddImport("fmt")
		fmt.Fprintf(w, "if %s == nil {\n", name.Full())
		fmt.Fprintf(w, "	return fmt.Errorf(\"field %s is required, should not be nil\" )\n", name.FieldName())
		fmt.Fprintf(w, "}\n")
		t.InnerType.Generate(w, cfg, name.WithPointer())
	}
}

func (t typePointer) Validate() error {
	return t.InnerType.Validate()
}
