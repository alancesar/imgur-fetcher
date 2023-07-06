package status

import "errors"

var (
	ErrNotFound  = errors.New("not found")
	ErrBadStatus = errors.New("bad status")
)
