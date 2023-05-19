package database

import (
	"database/sql"
	"strconv"

	_ "github.com/lib/pq"
)

func Connect() (*sql.DB, error) {
	db, err := sql.Open("postgres", "postgres://postgres:123456@localhost/test?sslmode=disable")
	if err != nil {
		return nil, err
	}
	return db, nil
}

func GetSquareQuery(latFlo, longFlo, radiusFlo float64) string {
	return `
		SELECT DISTINCT id, ST_X(coordinates::geometry) AS lat, ST_Y(coordinates::geometry) AS long
		FROM "MY_TABLE"
		WHERE ST_X(coordinates::geometry) >= ` + strconv.FormatFloat(latFlo-radiusFlo, 'f', -1, 64) +
		` AND ST_X(coordinates::geometry) <= ` + strconv.FormatFloat(latFlo+radiusFlo, 'f', -1, 64) +
		` AND ST_Y(coordinates::geometry) >= ` + strconv.FormatFloat(longFlo-radiusFlo, 'f', -1, 64) +
		` AND ST_Y(coordinates::geometry) <= ` + strconv.FormatFloat(longFlo+radiusFlo, 'f', -1, 64)
}

func GetCircleQuery(latFlo, longFlo, radiusFlo float64) string {
	return `
		SELECT DISTINCT id, ST_X(coordinates::geometry) AS lat, ST_Y(coordinates::geometry) AS long
		FROM "MY_TABLE"
		WHERE SQRT(ABS((` + strconv.FormatFloat(latFlo, 'f', -1, 64) + ` - ST_X(coordinates::geometry)) * (` + strconv.FormatFloat(latFlo, 'f', -1, 64) + ` - ST_X(coordinates::geometry)))
		+ ABS((` + strconv.FormatFloat(longFlo, 'f', -1, 64) + ` - ST_Y(coordinates::geometry)) * (` + strconv.FormatFloat(longFlo, 'f', -1, 64) + ` - ST_Y(coordinates::geometry)))) <= ` + strconv.FormatFloat(radiusFlo, 'f', -1, 64)
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
