package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	utl "siteparser/utils"
	"strconv"
	"time"
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

//AllUsersEmailsList struct for JESOnify
type AllUsersEmailsList struct {
	Emails []string `json:"emails"`
}

//GetUserEmails Test and get data from postgresql DB
func GetUserEmails(w http.ResponseWriter, r *http.Request) {
	// db := utl.DB
	log.Println("ID:", &utl.DB)
	// if err != nil {
	// 	log.Println("No connection to db! Error:", err)
	// } else {
	// 	log.Println("Get DB connection!")
	// 	// log.Println("ID:", *db)
	// }
	emailsArray := make([]string, 0)
	limit, lerr := strconv.Atoi(r.FormValue("limit"))
	offset, oerr := strconv.Atoi(r.FormValue("offset"))
	if lerr != nil {
		limit = 20
	}
	if oerr != nil {
		offset = 0
	}
	log.Println("LIMIT", limit)
	log.Println("OFFSET", offset)

	// defer db.Close()
	query := `
		SELECT email
		FROM authenticate_customuser
		WHERE email !='' OR email !=null
		LIMIT $1 
		OFFSET $2 
	`

	data, err := utl.DB.Query(query, limit, offset)
	defer data.Close()

	if err != nil {
		log.Println("Db Query error: ", err)
		http.Error(w, "Can not get Db data!", 400)
	} else {

		for data.Next() {

			var em string
			if err := data.Scan(&em); err != nil {
				// Check for a scan error.
				// Query rows will be closed with defer.
				log.Println("Error:", err)
			}
			emailsArray = append(emailsArray, em)
		}
		emailsFullObj := AllUsersEmailsList{Emails: emailsArray}

		log.Println("JSON:", emailsFullObj)

		w.Header().Set("Content-Type", "application/json")

		json.NewEncoder(w).Encode(emailsFullObj)
	}
}

//StructCreateNewUser New use structure
type StructCreateNewUser struct {
	Password       string    `json:"password"`
	IsSuperuser    bool      `json:"is_superuser"`
	Email          string    `json:"email"`
	Phone          string    `json:"phone"`
	IsStaff        bool      `json:"is_staff"`
	IsActive       bool      `json:"is_active"`
	DateJoined     time.Time `json:"date_joined"`
	DateUpdate     time.Time `json:"date_update"`
	EmailConfirm   bool      `json:"email_confirm"`
	AgentsIsActive bool      `json:"agents_is_active"`
}

//DetailUserDataStruct detail user data
type DetailUserDataStruct struct {
	ID             int       `json:"id"`
	IsSuperuser    bool      `json:"is_superuser"`
	Email          string    `json:"email"`
	Phone          string    `json:"phone"`
	IsStaff        bool      `json:"is_staff"`
	IsActive       bool      `json:"is_active"`
	DateJoined     time.Time `json:"date_joined"`
	DateUpdate     time.Time `json:"date_update"`
	EmailConfirm   bool      `json:"email_confirm"`
	AgentsIsActive bool      `json:"agents_is_active"`
}

//GetDetailUsersData get detail data for all users
func GetDetailUsersData(w http.ResponseWriter, r *http.Request) {

	limit, lerr := strconv.Atoi(r.FormValue("limit"))
	offset, oerr := strconv.Atoi(r.FormValue("offset"))
	if lerr != nil {
		limit = 20
	}
	if oerr != nil {
		offset = 0
	}
	log.Println("LIMIT", limit)
	log.Println("OFFSET", offset)
	query := `
	SELECT (id, email, phone, is_active, date_joined, date_updated)
	FROM authenticate_customuser
	LIMIT $1 
	OFFSET $1 
	`

	data, err := utl.DB.Query(query, limit, offset)

	for data.Next() {

	}

	if err != nil {
		log.Println("Error", err)
	} else {
		w.Header().Set("Content-Type", "application/json")

		json.NewEncoder(w).Encode(data)
	}
}

//CreateNewUser Create a new user
// func CreateNewUser(w http.ResponseWriter, r *http.Request) {

// 	// READ data from request must be added!!!

// 	query := `
// 		INSERT INTO authenticate_customuser ()
// 		VALUES authenticate_customuser
// 	`
// 	data, err := utl.DB.Exeq(query)
// 	if err != nil {
// 		log.Println("Error", err)
// 	} else {
// 		log.Println("Log data", data)
// 	}

// }
