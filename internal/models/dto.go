package models

type CreateLinkRequest struct{
	Url string `json:"url"`
}

type CreateLinkResponse struct{
	ShortCode string
}