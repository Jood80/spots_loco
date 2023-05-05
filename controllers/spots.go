package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/lib/pq"
)

type Spot struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Website     sql.NullString  `json:"website"`
	Coordinates string  `json:"coordinates"`
	Description sql.NullString   `json:"description"`
	Rating      float64 `json:"rating"`
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

	db, err := sql.Open("postgres", "postgres://postgres:123456@localhost/test?sslmode=disable")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer db.Close()

	query := fmt.Sprintf(`SELECT * FROM public."MY_TABLE" limit 10;`)

	// Execute SQL statement
	rows, err := db.Query(query)
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok && pqErr.Code == "42P01" {
			http.Error(w, "MY_TABLE table does not exist", http.StatusNotFound)
			} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	defer rows.Close()

	// Parse results into slice of Spot objects
	spots := []Spot{}
	for rows.Next() {
		spot := Spot{}
		err := rows.Scan(&spot.ID, &spot.Name, &spot.Website, &spot.Coordinates, &spot.Description, &spot.Rating)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		spots = append(spots, spot)
	}

	// Convert slice of Spot objects to JSON and return it
	jsonData, err := json.Marshal(spots)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}
