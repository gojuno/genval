package types

import (
	"fmt"
	"io"
)

func NewStructType(typeName string) *typeStruct {
	return &typeStruct{typeName: typeName}
}

type typeStruct struct {
	typeName string
	Method   *string
}

func (t typeStruct) Type() string {
	return t.typeName
}

func (t *typeStruct) SetTag(tag Tag) error {
	switch tag.Key() {
	case StructMethodKey:
		st := tag.(SimpleTag)
		t.Method = &st.Param
	default:
		return ErrUnusedTag
	}
	return nil
}

func (t typeStruct) Generate(w io.Writer, cfg GenConfig, name Name) {
	switch {
	case t.Method != nil:
		fmt.Fprintf(w, "if err := %s.%s(); err != nil {\n", name.WithoutPointer(), *t.Method)
		fmt.Fprintf(w, "    return err\n")
		fmt.Fprintf(w, "}\n")
	case cfg.NeedValidatableCheck:
		fmt.Fprintf(w, "if err := callValidateIfValidatable(%s); err != nil {\n", name.WithoutPointer())
		fmt.Fprintf(w, "    return err\n")
		fmt.Fprintf(w, "}\n")
	default:
		fmt.Fprintf(w, "if err := %s.Validate(); err != nil {\n", name.WithoutPointer())
		fmt.Fprintf(w, "    return err\n")
		fmt.Fprintf(w, "}\n")
	}
}

func (t typeStruct) Validate() error {
	return nil
}
