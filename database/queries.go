package database

import (
	"database/sql"
	"fmt"

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

func GetSquareQuery(latFlo, longFlo, radiusFlo float64) string {
	return fmt.Sprintf(`
	SELECT DISTINCT id, ST_X(coordinates::geometry) AS lat, ST_Y(coordinates::geometry) AS long
	FROM "MY_TABLE"
	WHERE ST_X(coordinates::geometry) >= %f AND ST_X(coordinates::geometry) <= %f
	AND ST_Y(coordinates::geometry) >= %f AND ST_Y(coordinates::geometry) <= %f`,
		latFlo-radiusFlo, latFlo+radiusFlo, longFlo-radiusFlo, longFlo+radiusFlo)
}

func GetCircleQuery(latFlo, longFlo, radiusFlo float64) string {
	return fmt.Sprintf(`
	SELECT DISTINCT id, ST_X(coordinates::geometry) AS lat, ST_Y(coordinates::geometry) AS long
	FROM "MY_TABLE"
	WHERE SQRT(ABS((%f - ST_X(coordinates::geometry)) * (%f - ST_X(coordinates::geometry)))
	+ ABS((%f - ST_Y(coordinates::geometry)) * (%f - ST_Y(coordinates::geometry)))) <= %f`,
		latFlo, latFlo, longFlo, longFlo, radiusFlo)
}

func GetThirdQuery() string {
	return `
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
}

func ExecuteShapeQuery(shape string, db *sql.DB, latFlo, longFlo, radiusFlo float64) (*sql.Rows, error) {
	var query string

	switch shape {
	case "square":
		query = GetSquareQuery(latFlo, longFlo, radiusFlo)
	case "circle":
		query = GetCircleQuery(latFlo, longFlo, radiusFlo)
	default:
		return nil, fmt.Errorf("Invalid shape: %s", shape)
	}

	rows, err := db.Query(query)
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok && pqErr.Code == "42P01" {
			return nil, fmt.Errorf("MY_TABLE table does not exist")
		} else {
			return nil, err
		}
	}
	defer rows.Close()

	return rows, nil
}

func ExecuteThirdQuery(db *sql.DB, latFlo, longFlo, radiusFlo float64) ([]Spot, error) {
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

	rows, err := db.Query(thirdQuery, latFlo, longFlo, radiusFlo*radiusFlo)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var newData []Spot

	for rows.Next() {
		var data Spot

		err := rows.Scan(&data.ID, &data.Name, &data.Website, &data.Coordinates, &data.Rating, &data.Lat, &data.Long, &data.Distance)
		if err != nil {
			return nil, err
		}
		newData = append(newData, data)
	}

	return newData, nil
}
