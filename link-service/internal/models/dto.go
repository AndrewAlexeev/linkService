package models

import (
	"errors"
	"strings"
	"time"
)

type LinkDto struct {
	ShortCode string
	Url       string
	Visits    int32
	CreatedAt time.Time
}

type CreateLinkRequest struct {
	Url string `json:"url"`
}

func (r CreateLinkRequest) Validate() error {
	if strings.TrimSpace(r.Url) == "" {
		return errors.New("url is required")
	}

	return nil
}

type CreateLinkResponse struct {
	ShortCode string `json:"short_code"`
}

type LinkResponse struct {
	Url    string `json:"url"`
	Visits int32  `json:"visits"`
}

type LinkStatResponse struct {
	ShortCode string    `json:"short_code"`
	Url       string    `json:"url"`
	Visits    int32     `json:"visits"`
	CreatedAt time.Time `json:"created_at"`
}
