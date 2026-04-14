package main

import (
	"link-service/internal/config"
	"link-service/internal/controllers"
	"link-service/internal/repository"
	"link-service/internal/services"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

func main() {

	startUpHttpServer()

}

func startUpHttpServer() {
	dbconfig := config.InitDbConfig()

	repository, err := repository.Init(dbconfig)

	if err != nil {
		log.Fatal("Repositry init error: ", err)
	}

	var linkService = services.LinkService{
		LinkRepository: repository}

	var linkController = controllers.LinkController{
		LinkService: linkService}

	config := config.InitConfig()

	mux := http.NewServeMux()

	mux.HandleFunc("POST /links", linkController.HandleCreate)

	// {id} - это переменная пути
	mux.HandleFunc("GET /links/{short_code}", linkController.HandleFetchLinkByShortCode)
	log.Println("Start app on port: ", config.Port)
	var error = http.ListenAndServe(":"+config.Port, mux)
	if error != nil {
		log.Fatal("Server startup error: ", error)
	}
}
