package server

import (
	"link-service/internal/config"
	"link-service/internal/handlers"
	"link-service/internal/repository"
	"link-service/internal/services"
	"log"
	"net/http"
)

func StartUpHttpServer() {
	dbconfig := config.InitDbConfig()

	lr, err := repository.NewLinkRepository(dbconfig)

	if err != nil {
		log.Fatal("Repositry init error: ", err)
	}

	var ls = services.NewLinkService(lr)

	var lh = handlers.NewLinkHandler(ls)

	config := config.InitConfig()

	mux := http.NewServeMux()

	mux.HandleFunc("POST /links", lh.Create)

	mux.HandleFunc("GET /links/{short_code}", lh.GetLinkByShortCode)
	mux.HandleFunc("GET /links/{short_code}/stats", lh.GetStatsByShortCode)
	mux.HandleFunc("DELETE /links/{short_code}", lh.DeleteByShortCode)

	log.Println("Start app on port: ", config.Port)
	var error = http.ListenAndServe(":"+config.Port, mux)
	if error != nil {
		log.Fatal("Server startup error: ", error)
	}
}
