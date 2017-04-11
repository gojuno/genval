package types

import "io"

const Byte string = "byte"

func NewByte() *typeByte {
	return &typeByte{}
}

type typeByte struct {
}

func (t typeByte) Type() string {
	return Byte
}

func (t *typeByte) SetTag(tag Tag) error {
	return ErrUnusedTag
}

func (t typeByte) Generate(w io.Writer, cfg GenConfig, name Name) {

}

func (t typeByte) Validate() error {
	return nil
}
