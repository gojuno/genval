package complicated_without_check

import "fmt"

type Status string

const (
	StatusCreated Status = "created"
	StatusPending Status = "pending"
	StatusActive  Status = "active"
	StatusFailed  Status = "failed"
)

type State int

const (
	StateOk    State = 200
	StateError State = 400
)

func (s State) Validate() error { //overriding
	if s < StateOk {
		return nil
	}
	if s > StateOk && s < StateError {
		return nil
	}
	return fmt.Errorf("unrecognized state: %s", s)
}
