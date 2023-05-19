package main

import (
	"example/routes"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func main() {
	r := httprouter.New()
	routes.RegisterRoutes(r)

	log.Println("Server is listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
