package controllers

import (
	"encoding/json"
	"link-service/internal/models"
	"link-service/internal/services"
	"net/http"
)

type LinkController struct {
	LinkService services.LinkService
}

func (lc LinkController) HandleCreate(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
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
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

}
