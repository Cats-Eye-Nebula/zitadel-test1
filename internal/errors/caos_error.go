package errors

import (
	"fmt"
)

var _ Error = (*CaosError)(nil)

type CaosError struct {
	Parent  error
	Message string
	ID      string
}

func ThrowError(parent error, id, message string) error {
	return createCaosError(parent, id, message)
}

func createCaosError(parent error, id, message string) *CaosError {
	return &CaosError{
		Parent:  parent,
		ID:      id,
		Message: message,
	}
}

func (err *CaosError) Error() string {
	if err.Parent != nil {
		return fmt.Sprintf("ID=%s Message=%s Parent=(%v)", err.ID, err.Message, err.Parent)
	}
	return fmt.Sprintf("ID=%s Message=%s", err.ID, err.Message)
}

func (err *CaosError) Unwrap() error {
	return err.GetParent()
}

func (err *CaosError) GetParent() error {
	return err.Parent
}

func (err *CaosError) GetMessage() string {
	return err.Message
}

func (err *CaosError) GetID() string {
	return err.ID
}
