package main

import "io"

func NewTypeFunc() *typeFunc {
	return &typeFunc{}
}

type typeFunc struct {
}

func (t *typeFunc) SetTag(tag Tag) error {
	return ErrUnusedTag
}

func (t typeFunc) Generate(w io.Writer, cfg GenConfig, suffix, name string) {

}
