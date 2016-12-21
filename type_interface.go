package main

import (
	"fmt"
	"io"
)

func NewTypeInterface() *typeInterface {
	return &typeInterface{}
}

type typeInterface struct {
	Func *string
}

func (t *typeInterface) SetTag(tag Tag) error {
	switch tag.Key() {
	case InterfaceFuncKey:
		st := tag.(SimpleTag)
		t.Func = &st.Param
	default:
		return ErrUnusedTag
	}
	return nil
}

func (t typeInterface) Generate(w io.Writer, cfg GenConfig, suffix, name string) {
	if t.Func != nil {
		fmt.Fprintf(w, "if err:=%s(%s); err!=nil {\n", *t.Func, suffix+name)
		fmt.Fprintf(w, "    return err\n")
		fmt.Fprintf(w, "}\n")
	}
}
