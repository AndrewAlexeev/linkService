package main
import (
"net/http"
"link-service/internal/controllers"
"link-service/internal/config"
"log"
"database/sql"
"fmt"
    _ "github.com/lib/pq"
)

func main(){
	
	startUpHttpServer()

}

func startUpHttpServer(){
		config := config.InitConfig()



		psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
        "password=%s dbname=%s sslmode=disable",
        config.DbHost,  config.DbPort,  config.DbUser,  config.DbPassword, config.DbName)

		log.Println(psqlInfo)
		

    	db, err := sql.Open("postgres", psqlInfo)

		if err != nil {
         log.Fatal("failed to connect to database: %w", err)
		 return
    	}

		result, err := db.Exec("insert into links (original_url, short_code, created_at, visits) values ($1,$2,NOW(), 0)", 
        "1","2")
    	if err != nil{
			log.Fatal("error: ", err)
        	panic(err)
    	}

		if (result != nil){
			log.Println("ok")
		}
    
    	if err := db.Ping(); err != nil {
         log.Fatal("failed to ping database: %w", err)
		 return
    	}

		http.HandleFunc("/links",controllers.HandleCreate)
		log.Println("Start app on port: ", config.Port)
		var error = http.ListenAndServe(":"+config.Port, nil)
		if error != nil {
        	log.Fatal("Server startup error: ", error)
    	}
}