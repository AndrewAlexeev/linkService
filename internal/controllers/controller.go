package controllers

import (
	"fmt"
	"net/http"
	"les1/internal/models"
	"les1/internal/services"
	"encoding/json"
)

func HandleCreate(w http.ResponseWriter, r *http.Request){

	if r.Method == http.MethodPost {
		var createLinkRequest models.CreateLinkRequest
		var err = json.NewDecoder(r.Body).Decode(&createLinkRequest)
		if err != nil {
			http.Error(w, "Error while parsing request", http.StatusBadRequest)
		}
	var createLinkResponse models.CreateLinkResponse
//TODO размер случайной ссылки берем из конфигов
	createLinkResponse.ShortCode = services.CreateShortLink(createLinkRequest.Url, 10)
	fmt.Fprintln(w, createLinkResponse)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

}