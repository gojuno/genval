package main

import (
	"fmt"
	"io"
)

func NewTypePointer(inner TypeDef) *typePointer {
	return &typePointer{
		Nullable:  false,
		InnerType: inner,
	}
}

type typePointer struct {
	Nullable  bool
	InnerType TypeDef
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

func (t typePointer) Generate(w io.Writer, cfg GenConfig, suffix, name string) {

	if t.Nullable {
		fmt.Fprintf(w, "if %s != nil {\n", suffix+name)
		t.InnerType.Generate(w, cfg, "*"+suffix, name)
		fmt.Fprintf(w, "}\n")
	} else {
		cfg.AddImport("fmt")
		fmt.Fprintf(w, "if %s == nil {\n", suffix+name)
		fmt.Fprintf(w, "	return fmt.Errorf(\"field %s is required, should not be nil\" )\n", name)
		fmt.Fprintf(w, "}\n")
		t.InnerType.Generate(w, cfg, "*"+suffix, name)
	}
}
