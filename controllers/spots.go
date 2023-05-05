package controllers

import (
	"encoding/json"
	"net/http"
)


func GetSpots(w http.ResponseWriter, r *http.Request) {
	lat := r.URL.Query().Get("lat")
	long := r.URL.Query().Get("long")
	radius := r.URL.Query().Get("radius")
	shape := r.URL.Query().Get("type")

	if lat == "" || long == "" || radius == "" || shape == "" {
		http.Error(w, "Missing query parameters", http.StatusBadRequest)
		return
	}


	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(lat+long+radius+shape)
}
