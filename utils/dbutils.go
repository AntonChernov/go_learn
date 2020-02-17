package utils

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
)

var (
	DB *sql.DB // Global variable handle connection to db
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

	if err = db.Ping(); err != nil {
		log.Panic(err)
	}

	return db, err

}

func init() {
	DB, _ = PostgressDBConnection()
}
