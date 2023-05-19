package controllers

import (
	"net/http"

	"example/service"

	"github.com/julienschmidt/httprouter"
)

func GetSpots(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	lat := r.Context().Value("lat").(float64)
	long := r.Context().Value("long").(float64)
	radius := r.Context().Value("radius").(float64)
	shape := r.Context().Value("shape").(string)

	jsonData, err := service.GetSpotData(lat, long, radius, shape)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}
