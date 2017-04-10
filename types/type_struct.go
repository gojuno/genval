package types

import (
	"fmt"
	"io"
	"strings"
)

func NewStruct(typeName string) *typeStruct {
	return &typeStruct{typeName: typeName, external: false}
}
func NewExternalStruct(typeName string) *typeStruct {
	return &typeStruct{typeName: typeName, external: true}
}

type typeStruct struct {
	typeName string
	external bool
	funcs    []string
}

func (t typeStruct) Type() string {
	return t.typeName
}

func (t *typeStruct) SetTag(tag Tag) error {
	switch tag.Key() {
	case StructFuncKey:
		for _, v := range parseFuncsParam(tag.(SimpleTag).Param) {
			t.funcs = append(t.funcs, v)
		}
	default:
		return ErrUnusedTag
	}
	return nil
}

func (t typeStruct) Generate(w io.Writer, cfg GenConfig, name Name) {
	switch {
	case len(t.funcs) != 0:
		for _, f := range t.funcs {
			if strings.HasPrefix(f, ".") {
				fmt.Fprintf(w, "if err := %s%s(); err != nil {\n", name.WithoutPointer(), f)
				fmt.Fprintf(w, "    return err\n")
				fmt.Fprintf(w, "}\n")
			} else {
				fmt.Fprintf(w, "if err := %s(%s); err != nil {\n", f, name.Full())
				fmt.Fprintf(w, "    return err\n")
				fmt.Fprintf(w, "}\n")
			}
		}
	case !cfg.NeedValidatableCheck, !t.external:
		fmt.Fprintf(w, "if err := %s.Validate(); err != nil {\n", name.WithAlias())
		fmt.Fprintf(w, "    return err\n")
		fmt.Fprintf(w, "}\n")
	default:
		fmt.Fprintf(w, "if err := validate(%s); err != nil {\n", name.WithAlias())
		fmt.Fprintf(w, "    return err\n")
		fmt.Fprintf(w, "}\n")
	}
}

func (t typeStruct) Validate() error {
	return nil
}
