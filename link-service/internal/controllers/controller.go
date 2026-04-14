package controllers

import (
	"encoding/json"
	"link-service/internal/models"
	"link-service/internal/services"
	"log"
	"net/http"
)

type LinkController struct {
	LinkService services.LinkService
}

func (lc LinkController) HandleCreate(w http.ResponseWriter, r *http.Request) {

	var createLinkRequest models.CreateLinkRequest
	var err = json.NewDecoder(r.Body).Decode(&createLinkRequest)
	if err != nil {
		http.Error(w, "Error while parsing request", http.StatusBadRequest)
		return
	}

	err = createLinkRequest.Validate()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	shortCode, err := lc.LinkService.SaveUrl(createLinkRequest.Url)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var createLinkResponse models.CreateLinkResponse
	createLinkResponse.ShortCode = shortCode
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(createLinkResponse)

}

func (lc LinkController) HandleFetchLinkByShortCode(w http.ResponseWriter, r *http.Request) {
	shortCode := r.PathValue("short_code")

	if shortCode == "" {
		http.Error(w, "short_code is nil", http.StatusBadRequest)
		return
	}

	linkDto, err := lc.LinkService.FindLinkByShortCode(shortCode)

	if err != nil {
		log.Println("Eror while fetch link by short code: %s error: %s", shortCode, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var linkResponse models.LinkResponse
	linkResponse.Url = linkDto.Url
	linkResponse.Visits = linkDto.Visits
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(linkResponse)
}
