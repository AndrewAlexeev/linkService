package main
import (
"net/http"
"les1/internal/controllers"
"log"

)

func main(){
	
	startUpHttpServer()

}

func startUpHttpServer(){
		http.HandleFunc("/links",controllers.HandleCreate)

		var error = http.ListenAndServe(":80", nil)
		if error != nil {
        	log.Fatal("Server startup error:", error)
    	}
}