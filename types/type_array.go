package types

import (
	"fmt"
	"go/ast"
	"io"
)

const Array string = "array"

func NewArray(inner TypeDef) *typeArray {
	return &typeArray{innerType: inner}
}

type typeArray struct {
	min       *string
	max       *string
	innerType TypeDef
}

func (t typeArray) Type() string {
	return Array
}

func (t *typeArray) SetValidateTag(tag ValidatableTag) error {
	switch tag.Key() {
	case ArrayMinItemsKey:
		st := tag.(SimpleTag)
		t.min = &st.Param
	case ArrayMaxItemsKey:
		st := tag.(SimpleTag)
		t.max = &st.Param
	case ArrayItemKey:
		scope := tag.(ScopeTag)
		for _, it := range scope.InnerTags {
			if err := t.innerType.SetValidateTag(it); err != nil {
				return fmt.Errorf("set item tags failed for %+v, err: %s", it, err)
			}
		}
	default:
		return ErrUnusedTag
	}
	return nil
}

func (t typeArray) Generate(w io.Writer, cfg GenConfig, name Name) {
	if t.min != nil {
		if *t.min != "0" {
			fmt.Fprintf(w, "if len(%s) < %s {\n", name.Full(), *t.min)
			fmt.Fprintf(w, "    errs.AddFieldf(%s, \"less items than %s\")\n", name.LabelName(), *t.min)
			fmt.Fprintf(w, "}\n")
		}
	}
	if t.max != nil {
		fmt.Fprintf(w, "if len(%s) > %s {\n", name.Full(), *t.max)
		fmt.Fprintf(w, "    errs.AddFieldf(%s, \"more items than %s\")\n", name.LabelName(), *t.max)
		fmt.Fprintf(w, "}\n")
	}

	if needGenerate(t.innerType) || validMaxMin(t.max, t.min) {
		fmt.Fprintf(w, "for i, x := range %s {\n", name.Full())
		fmt.Fprintf(w, "	_ = i \n")
		fmt.Fprintf(w, "	_ = x \n")
		cfg.AddImport("fmt")
		t.innerType.Generate(w, cfg, NewIndexedName(name.labelName[1:len(name.labelName)-1], "i", "x", name.tagName, false))
		fmt.Fprintf(w, "}\n")
	}
}

func (t typeArray) Validate() error {
	if err := validateMinMax(
		t.min,
		t.max,
		func(min float64) error {
			if min < 0 {
				return fmt.Errorf("min items can't be less than 0: %f", min)
			}
			return nil
		},
		func(max float64) error {
			if max < 0 {
				return fmt.Errorf("max items can't be less than 0: %f", max)
			}
			return nil
		},
	); err != nil {
		return err
	}
	return t.innerType.Validate()
}

func (t typeArray) Expr() ast.Expr {
	return nil
}
