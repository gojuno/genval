package api

type Status string

const (
	StatusCreated Status = "created"
	StatusPending Status = "pending"
	StatusActive  Status = "active"
	StatusFailed  Status = "failed"
)
