package error

import (
	"fmt"
	"strings"
)

type LabelError struct {
	Label string
	Err   error
}

const (
	TypeConflict     = "CONFLICT"
	TypeNotFound     = "NOT FOUND"
	TypeUnauthorized = "UNAUTHORIZED"
	TypeWrongFormat  = "WRONG_FORMAT"
	TypeCreated      = "CREATED"
)

func (le *LabelError) Error() string {
	return fmt.Sprintf("[%s] %v", le.Label, le.Err)
}

func NewLabelError(label string, err error) error {
	return &LabelError{
		Label: strings.ToUpper(label),
		Err:   err,
	}
}
