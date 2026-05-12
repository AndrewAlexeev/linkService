package errs

import (
	"fmt"
)

type LinkServiceError struct {
	Code    int
	Message string
}

func (e *LinkServiceError) Error() string {
	return fmt.Sprintf("Code %d: %s", e.Code, e.Message)
}

func NotFoundUrlError() error {
	return &LinkServiceError{Code: 404, Message: "Not found url in db"}
}
