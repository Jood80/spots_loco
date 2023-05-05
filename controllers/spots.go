package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type TestResponse struct {
	Message string `json:"message"`
}

func Test(w http.ResponseWriter, r *http.Request) {
	response := TestResponse{
		Message: "This is a test endpoint",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}



func bodyHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	
	if err !=nil {
		http.Error(w, "Error reading body", http.StatusInternalServerError)
		return
	}

	bodyString := string(body)

	fmt.Println("Request body: ", bodyString)
	fmt.Fprintf(w, "You sent this in the request body: %s", bodyString)
}
