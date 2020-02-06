package api

import (
	"fmt"
	"log"
	"net/http"
	utl "siteparser/utils"
)

//HelloHandler just example
func HelloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s !", r.URL.Path[0:])
	// return
}

//TestRequestHandler just example
func TestRequestHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This is test!")
}

//GetUserEmails Test and get data from postgresql DB
func GetUserEmails(w http.ResponseWriter, r *http.Request) {
	db, err := utl.PostgressDBConnection()
	defer db.Close()
	emails := make([]string, 0)
	data, err := db.Query("SELECT DISTINCT(email) FROM authenticate_customuser ")

	for data.Next() {
		var em string
		if err := data.Scan(&em); err != nil {
			// Check for a scan error.
			// Query rows will be closed with defer.
			log.Println("Error:", err)
		}
		emails = append(emails, em)
	}

	if err != nil {
		log.Println("No connection to db! Error:", err)
	} else {
		log.Println("Get DB connection!")
	}

	if err != nil {
		log.Println("Db Query error: ", err)
		http.Error(w, "Can not get Db data!", 400)
	} else {
		fmt.Fprintf(w, "Data from db: %s", emails)
	}
}
