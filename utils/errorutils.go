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

func JSONError(w http.ResponseWriter, errormsg string, stauscode int, sucess bool) {
	errorList := make([]string, 0)

	errorList = append(errorList, errormsg)
	var res interface{}

	w.Header().Set("Content-Type", "application/json")

	switch {
	case sucess != true && stauscode != 0:
		res = JSONErrorStruct{errorList}
		w.WriteHeader(stauscode)
		log.Println(errormsg)
	// Next case it is for future now this case not usable
	case sucess != true && stauscode == 0:
		res = JSONErrorStruct{errorList}
		w.WriteHeader(400)
		log.Println(errormsg)
	default:
		w.WriteHeader(200)
		log.Println(errormsg)
		res = JSONSucessStruct{errorList}
	}

	json.NewEncoder(w).Encode(&res)
}
