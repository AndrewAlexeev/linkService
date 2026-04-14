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

	http.HandleFunc("/links", linkController.HandleCreate)
	log.Println("Start app on port: ", config.Port)
	var error = http.ListenAndServe(":"+config.Port, nil)
	if error != nil {
		log.Fatal("Server startup error: ", error)
	}
}
