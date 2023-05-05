package main

import (
	"example/routes"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)


func main() {
	r :=mux.NewRouter()
	routes.RegisterRoutes(r)
	
	log.Println("Server is listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
