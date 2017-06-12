package types

import (
	"fmt"
	"io"
)

func NewNumber(typeName string) *typeNumber {
	return &typeNumber{typeName: typeName}
}

type typeNumber struct {
	typeName string
	min      *string
	max      *string
}

func (t typeNumber) Type() string {
	return t.typeName
}
func (t *typeNumber) SetTag(tag Tag) error {
	switch tag.Key() {
	case NumberMinKey:
		st := tag.(SimpleTag)
		t.min = &st.Param
	case NumberMaxKey:
		st := tag.(SimpleTag)
		t.max = &st.Param
	default:
		return ErrUnusedTag
	}
	return nil
}

func (t typeNumber) Generate(w io.Writer, cfg GenConfig, name Name) {
	if t.min != nil {
		fmt.Fprintf(w, "if %s < %s {\n", name.Full(), *t.min)
		fmt.Fprintf(w, "	errs.AddFieldf(%s, \"less than %s\")\n", name.LabelName(), *t.min)
		fmt.Fprintf(w, "}\n")
	}
	if t.max != nil {
		fmt.Fprintf(w, "if %s > %s {\n", name.Full(), *t.max)
		fmt.Fprintf(w, "	errs.AddFieldf(%s, \"more than %s\")\n", name.LabelName(), *t.max)
		fmt.Fprintf(w, "}\n")
	}
}

func (t typeNumber) Validate() error {
	return validateMinMax(
		t.min,
		t.max,
		func(min float64) error {
			return nil
		},
		func(max float64) error {
			return nil
		},
	)
}
