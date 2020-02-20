package utils

import (
	"encoding/json"
	_ "fmt"
	"log"
	"net/http"
)

//ErrorStruct struct for JSON errors response
type JSONErrorStruct struct {
	Errors []string `json:"errors"`
}

type JSONSucessStruct struct {
	Sucess []string `json:"sucess"`
}

func JSONError(w http.ResponseWriter, errormsg string, stuscode int, sucess bool) {
	errorList := make([]string, 0)

	errorList = append(errorList, errormsg)
	var res interface{}
	if sucess != true {
		res = JSONErrorStruct{errorList}
	} else {
		res = JSONSucessStruct{errorList}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(stuscode)
	log.Println(errormsg)
	json.NewEncoder(w).Encode(&res)
}
