package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"link-service/internal/errs"
	"link-service/internal/models"
	"link-service/internal/services"
	"log"
	"net/http"
	"strconv"
)

const maxLimit = 500

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
		writeJSONError(w, "Error while parsing request", http.StatusBadRequest)
		return
	}

	err = createLinkRequest.Validate()
	if err != nil {
		writeJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	shortCode, err := lh.linkService.SaveLink(r.Context(), createLinkRequest.Url)

	if err != nil {
		log.Printf("Error while creating link: %s", err)
		writeJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var createLinkResponse models.CreateLinkResponse
	createLinkResponse.ShortCode = shortCode

	writeJSONResponse(w, createLinkResponse, http.StatusCreated)

	log.Printf("Create link by short code: %s", shortCode)

}

func (lh *LinkHandler) GetLinkByShortCode(w http.ResponseWriter, r *http.Request) {
	shortCode := r.PathValue("short_code")

	if shortCode == "" {
		writeJSONError(w, "short_code parameter is required", http.StatusBadRequest)
		return
	}

	linkDto, err := lh.linkService.FindLinkByShortCode(r.Context(), shortCode)

	if err != nil {
		log.Printf("Error while fetching link by short code: %s error: %s", shortCode, err)

		var notFoundErr *errs.NotFoundLinkError

		if errors.As(err, &notFoundErr) {
			writeJSONError(w, err.Error(), http.StatusNotFound)
			return
		}

		writeJSONError(w, err.Error(), http.StatusInternalServerError)

		return
	}

	var linkResponse models.LinkResponse
	linkResponse.Url = linkDto.Url
	linkResponse.Visits = linkDto.Visits

	writeJSONResponse(w, linkResponse, http.StatusOK)

}

func (lh *LinkHandler) GetStatsByShortCode(w http.ResponseWriter, r *http.Request) {
	shortCode := r.PathValue("short_code")
	if shortCode == "" {
		writeJSONError(w, "short_code parameter is required", http.StatusBadRequest)
		return
	}

	linkDto, err := lh.linkService.FindLinkStatsByShortCode(r.Context(), shortCode)

	if err != nil {

		log.Printf("Error while fetch link by short code: %s error: %s", shortCode, err)

		var notFoundErr *errs.NotFoundLinkError

		if errors.As(err, &notFoundErr) {
			writeJSONError(w, err.Error(), http.StatusNotFound)
			return
		}

		writeJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := models.LinkStatResponse{
		ShortCode: linkDto.ShortCode,
		Url:       linkDto.Url,
		CreatedAt: linkDto.CreatedAt,
		Visits:    linkDto.Visits,
	}

	writeJSONResponse(w, res, http.StatusOK)

}

func (lh *LinkHandler) DeleteByShortCode(w http.ResponseWriter, r *http.Request) {
	shortCode := r.PathValue("short_code")
	if shortCode == "" {
		writeJSONError(w, "short_code parameter is required", http.StatusBadRequest)
		return
	}

	err := lh.linkService.DeleteByShortCode(r.Context(), shortCode)

	if err != nil {
		log.Printf("Error while delete link by short code: %s error: %s", shortCode, err)
		writeJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("Delete link by short code: %s", shortCode)
	w.WriteHeader(http.StatusNoContent)

}

func (lh *LinkHandler) GetLinks(w http.ResponseWriter, r *http.Request) {

	params := r.URL.Query()

	limit, offset, err := parsePaginationParams(params.Get("limit"), params.Get("offset"))
	if err != nil {
		writeJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	linkDtos, err := lh.linkService.GetByPage(r.Context(), limit, offset)

	if err != nil {
		log.Printf("Error while getting links: %s", err)
		writeJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSONResponse(w, linkDtos, http.StatusOK)

}

func parsePaginationParams(limitStr, offsetStr string) (int, int, error) {
	limit := 10 // default
	if limitStr != "" {
		parsed, err := strconv.Atoi(limitStr)
		if err != nil {
			return 0, 0, fmt.Errorf("invalid limit parameter: must be a number")
		}
		if parsed < 1 {
			return 0, 0, fmt.Errorf("limit must be greater than 0")
		}
		if parsed > maxLimit {
			return 0, 0, fmt.Errorf("limit must not exceed %d", maxLimit)
		}
		limit = parsed
	}

	offset := 0
	if offsetStr != "" {
		parsed, err := strconv.Atoi(offsetStr)
		if err != nil {
			return 0, 0, fmt.Errorf("invalid offset parameter: must be a number")
		}
		if parsed < 0 {
			return 0, 0, fmt.Errorf("offset must be 0 or greater")
		}
		offset = parsed
	}

	return limit, offset, nil
}

func writeJSONError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func writeJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Failed to encode response: %v", err)
	}
}
