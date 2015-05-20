// Package msgboard provides a message board application for users to
// create and edit HTML pages shown on a screen.
package msgboard

import (
	"net/http"

	"github.com/gorilla/mux"
)

func init() {
	r := mux.NewRouter()
	r.HandleFunc("/api/pages", ListPages).Methods("GET")
	r.HandleFunc("/api/pages", CreatePage).Methods("POST")
	r.HandleFunc("/api/pages/{id}", GetPage).Methods("GET")
	r.HandleFunc("/api/pages/{id}", UpdatePage).Methods("PATCH")
	http.Handle("/api/", r)
}
