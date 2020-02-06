package utils

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

//PostgressDBConnection Connection to Postgress DB
func PostgressDBConnection() (*sql.DB, error) {
	connStr := "postgres://sparseruser:superunsequrepasswordneedchange@localhost/sparser?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Println("Database connection not created!")
		log.Fatal(err)
		return nil, err
	}

	return db, err

}
