package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func testHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	
	if err !=nil {
		http.Error(w, "Error reading body", http.StatusInternalServerError)
		return
	}

	bodyString := string(body)

	fmt.Println("Request body: ", bodyString)
	fmt.Fprintf(w, "You sent this in the request body: %s", bodyString)
}

func main() {
	homeHandler := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome ya ana!")
	}

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/test", testHandler)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
