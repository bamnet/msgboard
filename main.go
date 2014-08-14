// Package msgboard provides a message board application for users to
// create and edit HTML pages shown on a screen.
package msgboard

import (
	"net/http"

	"github.com/gorilla/mux"
)

func init() {
	r := mux.NewRouter()
	r.HandleFunc("/pages", ListPages).Methods("GET")
	r.HandleFunc("/pages", CreatePage).Methods("POST")
	r.HandleFunc("/pages/{id}", GetPage).Methods("GET")
	r.HandleFunc("/pages/{id}", UpdatePage).Methods("PATCH")
	http.Handle("/", r)
}
