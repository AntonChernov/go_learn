package api

import (
	"encoding/json"
	"fmt"
	utl "go_learn/utils"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
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
	log.Println("ID:", &utl.DB)
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
		utl.JSONError(w, "Can not get Db data!", 400, false)
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
	Password       string `json:"password"`
	IsSuperuser    bool   `json:"is_superuser"`
	Email          string `json:"email"`
	Phone          string `json:"phone"`
	IsStaff        bool   `json:"is_staff"`
	IsActive       bool   `json:"is_active"`
	DateJoined     string `json:"date_joined"`
	DateUpdate     string `json:"date_update"`
	EmailConfirm   bool   `json:"email_confirm"`
	AgentsIsActive bool   `json:"agents_is_active"`
	LastLogin      string `json:"last_login"`
}

//DetailUserDataStruct detail user data
type DetailUserDataJSONStruct struct {
	ID             uint8     `json:"id"`
	IsSuperuser    bool      `json:"is_superuser"`
	Email          string    `json:"email"`
	Phone          string    `json:"phone"`
	IsStaff        bool      `json:"is_staff"`
	IsActive       bool      `json:"is_active"`
	DateJoined     time.Time `json:"date_joined"`
	DateUpdate     time.Time `json:"date_update"`
	EmailConfirm   bool      `json:"email_confirm"`
	AgentsIsActive bool      `json:"agents_is_active"`
	LastLogin      time.Time `json:"last_login"`
}

//GetDetailUsersData get detail data for all users
func GetDetailUsersData(w http.ResponseWriter, r *http.Request) {
	// Need create a Middlevar or something licke this bad way to use code duplication!!!
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
	SELECT id, email, phone, is_active, date_joined, date_update, is_staff, email_confirm, agents_is_active, last_login
	FROM authenticate_customuser
	ORDER BY 
	   id DESC
	LIMIT $1 
	OFFSET $2 
	`
	usersArray := make([]*DetailUserDataJSONStruct, 0)

	data, err := utl.DB.Query(query, limit, offset)
	defer data.Close()
	if err != nil {
		log.Println("DB connection lost!", err)
		utl.JSONError(w, "DB connection lost!", 400, false)
	}
	//List of users use deklared structure
	for data.Next() {
		obj := new(DetailUserDataJSONStruct)
		err := data.Scan(&obj.ID, &obj.Email, &obj.Phone, &obj.IsActive, &obj.DateJoined, &obj.DateUpdate, &obj.IsStaff, &obj.EmailConfirm, &obj.AgentsIsActive, &obj.LastLogin)

		if err != nil {
			log.Println("Error", err)
			utl.JSONError(w, "Can not get data from DB!", 400, false)
			return
		}
		usersArray = append(usersArray, obj)
	}
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(usersArray)
}

//CreateNewUser Create a new user
func CreateNewUser(w http.ResponseWriter, r *http.Request) {

	query := `
		INSERT INTO authenticate_customuser(email, phone, password, is_active, date_joined, date_update, is_staff, email_confirm, agents_is_active, last_login, is_superuser)
		VALUES ($1,$2,$3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, email, phone, is_active, date_joined, date_update, is_staff, email_confirm, agents_is_active, last_login
	`

	decoder := json.NewDecoder(r.Body)
	var user StructCreateNewUser
	err := decoder.Decode(&user)
	obj := new(DetailUserDataJSONStruct) //user object
	if err != nil {
		panic(err)
	}
	log.Println(&user)

	data := utl.DB.QueryRow(
		query, &user.Email, &user.Phone, &user.Password, &user.IsActive, time.Now().Format(time.RFC3339),
		time.Now().Format(time.RFC3339), &user.IsStaff, &user.EmailConfirm, &user.AgentsIsActive,
		time.Now().Format(time.RFC3339), &user.IsSuperuser).Scan(&obj.ID, &obj.Email, &obj.Phone, &obj.IsActive,
		&obj.DateJoined, &obj.DateUpdate, &user.IsStaff, &user.EmailConfirm, &obj.AgentsIsActive, &obj.LastLogin) //create user obj and return fields related to this instance

	if data != nil {
		log.Println("Error", err)
	} else {
		log.Println("Log data", &obj)
		w.Header().Set("Content-Type", "application/json")

		json.NewEncoder(w).Encode(&obj)
	}

}

//UpdatePhoneEmailStruct struct for update email and phone for some user
type UpdatePhoneEmailStruct struct {
	Phone string `json:"phone"`
	Email string `json:"email"`
}

func UpdateUserView(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userId := params["id"]

	decoder := json.NewDecoder(r.Body)
	var user UpdatePhoneEmailStruct
	err := decoder.Decode(&user)

	query := `
		UPDATE authenticate_customuser
		SET phone = $2, email = $3
		WHERE id=$1
		RETURNING id, email, phone, is_active, date_joined, date_update, is_staff, email_confirm, agents_is_active, last_login
	`
	// log.Println(&user)
	data, err := utl.DB.Query(query, &userId, &user.Phone, &user.Email)
	defer data.Close()
	if err != nil {
		log.Panic(err)
		// log.Error("Error", err)
		utl.JSONError(w, "DB error", 400, false)
	} else {
		obj := new(DetailUserDataJSONStruct)
		for data.Next() {

			err := data.Scan(&obj.ID, &obj.Email, &obj.Phone, &obj.IsActive, &obj.DateJoined, &obj.DateUpdate, &obj.IsStaff, &obj.EmailConfirm, &obj.AgentsIsActive, &obj.LastLogin)

			if err != nil {
				log.Println("Error", err)
				utl.JSONError(w, "Can not get data from DB!", 400, false)
				return
			}

		}
		w.Header().Set("Content-Type", "application/json")

		json.NewEncoder(w).Encode(&obj)
	}

}

func DeleteUserView(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userId := params["id"]
	log.Printf("Type ID %T", userId)
	// queryCheckUserExist := `
	// SELECT
	// EXISTS (SELECT 1 FROM authenticate_customuser WHERE id = $1)
	// `
	queryCheckUserExist := `
	SELECT id 
	FROM authenticate_customuser 
	WHERE id = $1
	`
	var id int
	err := utl.DB.QueryRow(queryCheckUserExist, userId).Scan(&id)
	log.Println(err)

	if err != nil {
		log.Println(err)
		utl.JSONError(w, "User Does Not Exist!", 400, false) // TODO Create own error for return JSON with Error text for front end devs!!! DONE utils/errorutils
		// return
	} else {
		query := `
		DELETE FROM authenticate_customuser
		WHERE id = $1
		`
		_, derr := utl.DB.Query(query, userId) // finde a way how to resolve not exist user deletion!!!
		if derr != nil {
			log.Panic(derr)
		} else {
			utl.JSONError(w, "User succesfulli deleted!", 200, true)
			// return
		}
	}

}
