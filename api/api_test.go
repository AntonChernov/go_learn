package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"

	utl "go_learn/utils"
)

func TestIndexHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/", HelloHandler)
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

}

func TestTestHandler(t *testing.T) {

	url := "/test/"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc(url, TestRequestHandler)
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	returnedInBody, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		t.Errorf("Error when do get requert to %v", req.URL)
	}
	example := "This is test!"
	respBpbyString := string(returnedInBody)
	if respBpbyString != example {
		t.Errorf("Body incorrect! Response boby %v - Must be \"This is test!\" ", respBpbyString)
	}
}

func TestDBHandler(t *testing.T) {
	url := "/test-db/"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc(url, GetUserEmails)
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestBadDBHandler(t *testing.T) {
	url := "/test-db/"
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc(url, GetUserEmails).Methods("GET")
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusMethodNotAllowed)
	}
}

func TestUsersListHandler(t *testing.T) {
	url := "/users-list/"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc(url, GetUserEmails)
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestBadRequestUsersListHandler(t *testing.T) {
	url := "/users-list/"
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc(url, GetUserEmails).Methods("GET")
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusMethodNotAllowed)
	}
}

func TestUsersCreateDeleteHandler(t *testing.T) {
	// user created success
	url := "/create-user/"
	rand.Seed(time.Now().UnixNano())
	emailStrGenerate := fmt.Sprintf("som1e_valid%vEmail@gmail.com", rand.Intn(10000))
	userData := StructCreateNewUser{Email: emailStrGenerate, Phone: "56468768767789", Password: "12345", IsActive: true, DateJoined: "2020-02-20T09:20:30Z",
		DateUpdate: "2020-02-20T09:20:30Z", IsStaff: false, EmailConfirm: false, AgentsIsActive: false, IsSuperuser: false}
	body, mErr := json.Marshal(userData)
	if mErr != nil {
		t.Fatal(mErr)
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc(url, CreateNewUser)
	router.ServeHTTP(rr, req)

	decoder := json.NewDecoder(rr.Body)
	var user DetailUserDataJSONStruct
	eerr := decoder.Decode(&user)

	if eerr != nil {
		t.Fatal(eerr)
	}
	if user.Email != emailStrGenerate {
		t.Errorf("User email missmatch! Shoulbe: %v but reade from json %v", emailStrGenerate, user.Email)
	}

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	//delete user success
	delURL := fmt.Sprintf("/delete-user/%v", user.ID)
	fmt.Println("URL", delURL)
	dReq, dErr := http.NewRequest("POST", delURL, nil)
	if dErr != nil {
		t.Error(dErr)
	}

	dRr := httptest.NewRecorder()

	dRouter := mux.NewRouter()
	dRouter.HandleFunc("/delete-user/{id}", DeleteUserView)
	fmt.Println("Router URL", dReq.RequestURI)
	dRouter.ServeHTTP(dRr, dReq)

	if dStatus := dRr.Code; dStatus != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			dStatus, http.StatusOK)
	}

	dDecoder := json.NewDecoder(dRr.Body)
	var dUser utl.JSONSucessStruct
	duErr := dDecoder.Decode(&dUser)

	if duErr != nil {
		t.Error(duErr)
	}

	goodMsg := "User succesfulli deleted!"

	if dUser.Sucess[0] != goodMsg {
		t.Errorf("Return strings missmatch! Shoud be %v but read from JSON %v", goodMsg, dUser.Sucess[0])
	}

	// del return error because user dont exist
	delErrURL := fmt.Sprintf("/delete-user/%v", user.ID)

	delErrReq, delErr := http.NewRequest("POST", delErrURL, nil)
	if delErr != nil {
		t.Error(delErr)
	}

	delErrRr := httptest.NewRecorder()

	delErrRouter := mux.NewRouter()
	delErrRouter.HandleFunc("/delete-user/{id}", DeleteUserView)
	delErrRouter.ServeHTTP(delErrRr, delErrReq)

	if dErrStatus := delErrRr.Code; dErrStatus != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			dErrStatus, http.StatusBadRequest)
	}

	dErrDecoder := json.NewDecoder(delErrRr.Body)
	var dErrUser utl.JSONErrorStruct
	dUserErr := dErrDecoder.Decode(&dErrUser)

	if dUserErr != nil {
		t.Error(dUserErr)
	}

	badMsg := "User Does Not Exist!"

	if dErrUser.Errors[0] != badMsg {
		t.Errorf("Return strings missmatch! Shoud be %v but read from JSON %v", badMsg, dErrUser.Errors[0])
	}

}

func TestBadRequestUsersCreateHandler(t *testing.T) {
	// get error because body is empty
	url := "/create-user/"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc(url, GetUserEmails).Methods("POST")
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusMethodNotAllowed)
	}
}

func TestCreateUpdateDeleteUser(t *testing.T) {
	//code here
}
