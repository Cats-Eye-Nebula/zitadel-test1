package errors

import "fmt"

type NotFound interface {
	error
	IsNotFound()
}

type NotFoundError struct {
	*CaosError
}

func ThrowNotFound(parent error, id, message string) error {
	return &NotFoundError{CreateCaosError(parent, id, message)}
}

func ThrowNotFoundf(parent error, id, format string, a ...interface{}) error {
	return ThrowNotFound(parent, id, fmt.Sprintf(format, a...))
}

func (err *NotFoundError) IsNotFound() {}

func IsNotFound(err error) bool {
	_, ok := err.(NotFound)
	return ok
}
