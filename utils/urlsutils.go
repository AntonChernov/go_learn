package utils

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
)

//Handler struct for build chane from URL path and handler it self
type Handler struct {
	path    string
	f       http.HandlerFunc
	methods []string
	name    string
}

//MakeURLRoutersToHandlers So all detailed info in name ^)
func MakeURLRoutersToHandlers(hs []Handler, routinst *mux.Router) {
	for _, h := range hs {
		if len(h.methods) == 0 {
			routinst.HandleFunc(h.path, h.f).Name(h.name)
		} else {
			routinst.HandleFunc(h.path, h.f).Methods(h.methods...).Name(h.name)
		}
	}
}

//FilePathWalker for future when i will better know go
func FilePathWalker(pattern string) []string {

	var files []string
	var res []string

	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		_, fileName := filepath.Split(file)

		if fileName == pattern {
			res = append(res, file)
		}

	}
	return res
}
