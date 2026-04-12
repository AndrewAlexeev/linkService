package main
import (
"net/http"
"link-service/internal/controllers"
"log"

)

func main(){
	
	startUpHttpServer()

}

func startUpHttpServer(){
		http.HandleFunc("/links",controllers.HandleCreate)
		log.Println("Start app")
		var error = http.ListenAndServe(":8080", nil)
		if error != nil {
        	log.Fatal("Server startup error: ", error)
    	}
}