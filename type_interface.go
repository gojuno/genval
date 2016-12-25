package main

import (
	"fmt"
	"io"
)

func NewInterfaceType() *typeInterface {
	return &typeInterface{}
}

type typeInterface struct {
	Func *string
}

func (t typeInterface) Type() string {
	return "interface"
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

func (t typeInterface) Generate(w io.Writer, cfg GenConfig, name Name) {
	if t.Func != nil {
		fmt.Fprintf(w, "if err:=%s(%s); err!=nil {\n", *t.Func, name.Full())
		fmt.Fprintf(w, "    return err\n")
		fmt.Fprintf(w, "}\n")
	}
}

func (t typeInterface) Validate() error {
	return nil
}
