package errs

import (
	"fmt"
)

type NotFoundLinkError struct {
	Message string
}

func (e *NotFoundLinkError) Error() string {
	return fmt.Sprintf("error info: %s", e.Message)
}

func NewNotFoundLinkError() error {
	return &NotFoundLinkError{Message: "Not found url in db"}
}
