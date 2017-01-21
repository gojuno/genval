package types

import "io"

func NewChan() *typeChan {
	return &typeChan{}
}

type typeChan struct {
}

func (t typeChan) Type() string {
	return "chan"
}

func (t *typeChan) SetTag(tag Tag) error {
	return ErrUnusedTag
}

func (t typeChan) Generate(w io.Writer, cfg GenConfig, name Name) {

}

func (t typeChan) Validate() error {
	return nil
}
