package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/lib/pq"
)

type Spot struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Website     sql.NullString `json:"website"`
	Coordinates string         `json:"coordinates"`
	// Description sql.NullString `json:"description"`
	Rating   float64 `json:"rating"`
	Lat      float64 `json:"lat"`
	Long     float64 `json:"long"`
	Distance float64 `json:"distance"`
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

	db, err := sql.Open("postgres", "postgres://postgres:123456@localhost/test?sslmode=disable")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer db.Close()

	var query string
	switch shape {
	case "square":
		query = fmt.Sprintf(`
			SELECT DISTINCT id, ST_X(coordinates::geometry) AS lat, ST_Y(coordinates::geometry) AS long
			FROM "MY_TABLE"
			WHERE ST_X(coordinates::geometry) >= %f AND ST_X(coordinates::geometry) <= %f
			AND ST_Y(coordinates::geometry) >= %f AND ST_Y(coordinates::geometry) <= %f`,
			latFlo-radiusFlo, latFlo+radiusFlo, longFlo-radiusFlo, longFlo+radiusFlo)
	case "circle":
		query = fmt.Sprintf(`
			SELECT DISTINCT id, ST_X(coordinates::geometry) AS lat, ST_Y(coordinates::geometry) AS long
			FROM "MY_TABLE"
			WHERE SQRT(ABS((%f - ST_X(coordinates::geometry)) * (%f - ST_X(coordinates::geometry)))
			+ ABS((%f - ST_Y(coordinates::geometry)) * (%f - ST_Y(coordinates::geometry)))) <= %f`,
			latFlo, latFlo, longFlo, longFlo, radiusFlo)
	}
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

	thirdQuery := `
	SELECT
		id,
		name,
		website,
		coordinates,
		rating,
		ST_X(coordinates::geometry) AS lat,
		ST_Y(coordinates::geometry) AS long,
		SQRT(ABS(($1 - ST_X(coordinates::geometry)) * ($1 - ST_X(coordinates::geometry))) + ABS(($2 - ST_Y(coordinates::geometry)) * ($2 - ST_Y(coordinates::geometry)))) AS distance
	FROM public."MY_TABLE"
	WHERE SQRT(ABS(($1 - ST_X(coordinates::geometry)) * ($1 - ST_X(coordinates::geometry))) + ABS(($2 - ST_Y(coordinates::geometry)) * ($2 - ST_Y(coordinates::geometry)))) <= $3
	GROUP BY id, name, website, coordinates, rating, lat, long
	ORDER BY
		CASE
			WHEN SQRT(ABS(($1 - ST_X(coordinates::geometry)) * ($1 - ST_X(coordinates::geometry))) + ABS(($2 - ST_Y(coordinates::geometry)) * ($2 - ST_Y(coordinates::geometry)))) <= 50
			THEN 0
			ELSE SQRT(ABS(($1 - ST_X(coordinates::geometry)) * ($1 - ST_X(coordinates::geometry))) + ABS(($2 - ST_Y(coordinates::geometry)) * ($2 - ST_Y(coordinates::geometry))))
		END ASC, rating DESC`

	thirdRows, err := db.Query(thirdQuery, latFlo, longFlo, radiusFlo*radiusFlo)
	if err != nil {
		println("Error querying")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal(err)
	}
	defer thirdRows.Close()

	for thirdRows.Next() {
		println("we're here")
		var data Spot

		err := thirdRows.Scan(&data.ID, &data.Name, &data.Website, &data.Coordinates, &data.Rating, &data.Lat, &data.Long, &data.Distance)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		newData = append(newData, data)
		println("rowData", rowData)
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
