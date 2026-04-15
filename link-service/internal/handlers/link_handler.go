package handlers

import (
	"encoding/json"
	"link-service/internal/models"
	"link-service/internal/services"
	"log"
	"net/http"
)

type LinkHandler struct {
	linkService *services.LinkService
}

func NewLinkHandler(linkService *services.LinkService) LinkHandler {
	return LinkHandler{linkService: linkService}
}

func (lh LinkHandler) Create(w http.ResponseWriter, r *http.Request) {

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

	shortCode, err := lh.linkService.SaveUrl(createLinkRequest.Url)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var createLinkResponse models.CreateLinkResponse
	createLinkResponse.ShortCode = shortCode
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(createLinkResponse)
	log.Printf("Create link by short code: %s", shortCode)

}

func (lh LinkHandler) GetLinkByShortCode(w http.ResponseWriter, r *http.Request) {
	shortCode := r.PathValue("short_code")

	if shortCode == "" {
		http.Error(w, "short_code is nil", http.StatusBadRequest)
		return
	}

	linkDto, err := lh.linkService.FindLinkByShortCode(shortCode)

	if err != nil {
		log.Printf("Eror while fetch link by short code: %s error: %s", shortCode, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var linkResponse models.LinkResponse
	linkResponse.Url = linkDto.Url
	linkResponse.Visits = linkDto.Visits
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(linkResponse)
}

func (lh LinkHandler) GetStatsByShortCode(w http.ResponseWriter, r *http.Request) {
	shortCode := r.PathValue("short_code")
	if shortCode == "" {
		http.Error(w, "short_code is nil", http.StatusBadRequest)
		return
	}

	linkDto, err := lh.linkService.FindLinkStatsByShortCode(shortCode)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	res := models.LinkStatResponse{
		ShortCode: linkDto.ShortCode,
		Url:       linkDto.Url,
		CreatedAt: linkDto.CreatedAt,
		Visits:    linkDto.Visits,
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(res)

}

func (lh LinkHandler) DeleteByShortCode(w http.ResponseWriter, r *http.Request) {
	shortCode := r.PathValue("short_code")
	if shortCode == "" {
		http.Error(w, "short_code is nil", http.StatusBadRequest)
		return
	}

	err := lh.linkService.DeleteByShortCode(shortCode)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	log.Printf("Delete link by short code: %s", shortCode)

}
