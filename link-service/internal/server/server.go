package server

import (
	"link-service/internal/cache"
	"link-service/internal/config"
	"link-service/internal/handlers"
	"link-service/internal/repository"
	"link-service/internal/services"
	"log"
	"net/http"
)

func StartUpHttpServer() {

	redisConfig, errRedis := config.InitRedisConfig()

	if errRedis != nil {
		log.Fatal("Redis init error: ", errRedis)
	}

	lc := cache.InitLinkCache(*redisConfig)

	dbconfig, err := config.InitDbConfig()

	if err != nil {
		log.Fatal("db config init error: ", err)
	}

	lr, err := repository.NewLinkRepository(dbconfig)

	if err != nil {
		log.Fatal("Repositry init error: ", err)
	}

	var ls = services.NewLinkService(lr, lc)

	var lh = handlers.NewLinkHandler(ls)

	serverConfig := config.InitServerConfig()

	mux := http.NewServeMux()

	mux.HandleFunc("POST /links", lh.Create)

	mux.HandleFunc("GET /links/{short_code}", lh.GetLinkByShortCode)
	mux.HandleFunc("GET /links/{short_code}/stats", lh.GetStatsByShortCode)
	mux.HandleFunc("DELETE /links/{short_code}", lh.DeleteByShortCode)
	mux.HandleFunc("GET /links", lh.GetLinks)

	log.Println("Start app on port: ", serverConfig.Port)

	if err := http.ListenAndServe(":"+serverConfig.Port, mux); err != nil {
		log.Fatal("Server startup error: ", err)
	}
}
