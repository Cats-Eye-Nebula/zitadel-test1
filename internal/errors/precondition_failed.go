package errors

import (
	"fmt"
)

var (
	_ PreconditionFailed = (*PreconditionFailedError)(nil)
	_ Error              = (*PreconditionFailedError)(nil)
)

type PreconditionFailed interface {
	error
	IsPreconditionFailed()
}

type PreconditionFailedError struct {
	*CaosError
}

func ThrowPreconditionFailed(parent error, id, message string) error {
	return &PreconditionFailedError{CreateCaosError(parent, id, message)}
}

func ThrowPreconditionFailedf(parent error, id, format string, a ...interface{}) error {
	return ThrowPreconditionFailed(parent, id, fmt.Sprintf(format, a...))
}

func (err *PreconditionFailedError) IsPreconditionFailed() {}

func IsPreconditionFailed(err error) bool {
	_, ok := err.(PreconditionFailed)
	return ok
}
