package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func Hello(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	log.Println("Hello")
	fmt.Fprintf(w, "Server is running!")
}
