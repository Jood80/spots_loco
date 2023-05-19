package controllers

import (
	"net/http"
	"strconv"

	"example/service"
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

	latFlo, err := strconv.ParseFloat(lat, 64)
	longFlo, err := strconv.ParseFloat(long, 64)
	radiusFlo, err := strconv.ParseFloat(radius, 64)

	jsonData, err := service.GetSpotData(latFlo, longFlo, radiusFlo, shape)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}
