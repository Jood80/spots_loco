package controllers

import (
	"fmt"
	"log"
	"net/http"
)

func Hello(w http.ResponseWriter, r *http.Request) {
	log.Println("Hello")
	fmt.Fprintf(w, "Server is running!")
}
