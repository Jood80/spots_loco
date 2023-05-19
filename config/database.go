package config

import (
	"database/sql"
)

func ConnectDB() (*sql.DB, error) {
	db, err := sql.Open("postgres", "postgres://postgres:123456@localhost/test?sslmode=disable")
	if err != nil {
		return nil, err
	}
	return db, nil
}
