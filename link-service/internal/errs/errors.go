package errs

type NotFoundLinkError struct {
	Message string
}

func (e *NotFoundLinkError) Error() string {
	return e.Message
}

func NewNotFoundLinkError() error {
	return &NotFoundLinkError{Message: "Not found url in db"}
}
