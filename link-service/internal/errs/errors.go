package errs

import (
	"fmt"
)

type NotFoundLinkError struct {
	Code    int
	Message string
}

func (e *NotFoundLinkError) Error() string {
	return fmt.Sprintf("Code %d: %s", e.Code, e.Message)
}

func NotFoundUrlError() error {
	return &NotFoundLinkError{Code: 404, Message: "Not found url in db"}
}
