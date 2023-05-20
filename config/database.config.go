package config

import (
	"database/sql"
	"example/model"
	"log"
)

var schemaBuilt = false

func ConnectDB() (*sql.DB, error) {
	// Establish the database connection
	db, err := sql.Open("postgres", "postgres://postgres:123456@localhost/test2?sslmode=disable")
	if err != nil {
		return nil, err
	}

	// Ensure the connection is valid
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	// Build the database schema if it hasn't been built before
	if !schemaBuilt {
		err = model.BuildSchema(db)
		if err != nil {
			db.Close()
			return nil, err
		}
		schemaBuilt = true
		log.Println("Schema build completed successfully")
	}

	log.Println("Connection established")

	return db, nil
}
