package handlers

import (
	"encoding/json"
	"link-service/internal/errs"
	"link-service/internal/models"
	"link-service/internal/services"
	"log"
	"net/http"
)

type LinkHandler struct {
	linkService *services.LinkService
}

func NewLinkHandler(linkService *services.LinkService) *LinkHandler {
	return &LinkHandler{linkService: linkService}
}

func (lh *LinkHandler) Create(w http.ResponseWriter, r *http.Request) {

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

	shortCode, err := lh.linkService.SaveLink(r.Context(), createLinkRequest.Url)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var createLinkResponse models.CreateLinkResponse
	createLinkResponse.ShortCode = shortCode
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err = json.NewEncoder(w).Encode(createLinkResponse); err != nil {
		log.Printf("Failed to encode response: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	log.Printf("Create link by short code: %s", shortCode)

}

func (lh *LinkHandler) GetLinkByShortCode(w http.ResponseWriter, r *http.Request) {
	shortCode := r.PathValue("short_code")

	if shortCode == "" {
		http.Error(w, "short_code parameter is required", http.StatusBadRequest)
		return
	}

	linkDto, err := lh.linkService.FindLinkByShortCode(r.Context(), shortCode)

	if err != nil {
		log.Printf("Error while fetch link by short code: %s error: %s", shortCode, err)

		if _, ok := err.(*errs.NotFoundLinkError); ok {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var linkResponse models.LinkResponse
	linkResponse.Url = linkDto.Url
	linkResponse.Visits = linkDto.Visits
	w.Header().Set("Content-Type", "application/json")

	if err = json.NewEncoder(w).Encode(linkResponse); err != nil {
		log.Printf("Failed to encode response: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

}

func (lh *LinkHandler) GetStatsByShortCode(w http.ResponseWriter, r *http.Request) {
	shortCode := r.PathValue("short_code")
	if shortCode == "" {
		http.Error(w, "short_code parameter is required", http.StatusBadRequest)
		return
	}

	linkDto, err := lh.linkService.FindLinkStatsByShortCode(r.Context(), shortCode)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := models.LinkStatResponse{
		ShortCode: linkDto.ShortCode,
		Url:       linkDto.Url,
		CreatedAt: linkDto.CreatedAt,
		Visits:    linkDto.Visits,
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Printf("Failed to encode response: %v", err)
	}

}

func (lh *LinkHandler) DeleteByShortCode(w http.ResponseWriter, r *http.Request) {
	shortCode := r.PathValue("short_code")
	if shortCode == "" {
		http.Error(w, "short_code parameter is required", http.StatusBadRequest)
		return
	}

	err := lh.linkService.DeleteByShortCode(r.Context(), shortCode)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("Delete link by short code: %s", shortCode)

}

func (lh *LinkHandler) GetLinks(w http.ResponseWriter, r *http.Request) {

	params := r.URL.Query()

	limit := params.Get("limit")
	offset := params.Get("offset")

	linkDtos, err := lh.linkService.GetByPage(r.Context(), limit, offset)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(linkDtos); err != nil {
		log.Printf("Failed to encode response: %v", err)
	}

}
