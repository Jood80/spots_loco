package controllers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"example/database"

	"github.com/lib/pq"
)

type Spot struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Website     sql.NullString `json:"website"`
	Coordinates string         `json:"coordinates"`
	Rating      float64        `json:"rating"`
	Lat         float64        `json:"lat"`
	Long        float64        `json:"long"`
	Distance    float64        `json:"distance"`
}

type Rows struct {
	ID   string  `json:"id"`
	Lat  float64 `json:"lat"`
	Long float64 `json:"long"`
}

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

	db, err := database.Connect()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var query string
	switch shape {
	case "square":
		query = database.GetSquareQuery(latFlo, longFlo, radiusFlo)
	case "circle":
		query = database.GetCircleQuery(latFlo, longFlo, radiusFlo)
	}

	rows, err := db.Query(query)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "42P01" {
			http.Error(w, "MY_TABLE table does not exist", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	defer rows.Close()

	var rowData []Rows
	var newData []Spot

	for rows.Next() {
		var data Rows

		err := rows.Scan(&data.ID, &data.Lat, &data.Long)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Fatal(err)
		}

		rowData = append(rowData, data)
	}

	thirdQuery := database.GetThirdQuery()

	thirdRows, err := db.Query(thirdQuery, latFlo, longFlo, radiusFlo*radiusFlo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal(err)
	}
	defer thirdRows.Close()

	for thirdRows.Next() {
		var data Spot

		err := thirdRows.Scan(&data.ID, &data.Name, &data.Website, &data.Coordinates, &data.Rating, &data.Lat, &data.Long, &data.Distance)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		newData = append(newData, data)
	}

	if err = rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal(err)
	}

	jsonData, err := json.Marshal(newData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}
