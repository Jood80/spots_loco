package model

import (
	"database/sql"
	"io/ioutil"
	"log"
)

func BuildSchema(db *sql.DB) error {
	sqlFile, err := ioutil.ReadFile("model/spots.sql")
	if err != nil {
		log.Fatalf("Failed to read SQL file: %v", err)
	}
	// Execute the SQL statements
	_, err = db.Exec(string(sqlFile))
	if err != nil {
		return err
	}
	return nil
}
