package msgboard

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func init() {
	r := mux.NewRouter()
	r.HandleFunc("/pages", ListPages).Methods("GET")
	r.HandleFunc("/pages", CreatePage).Methods("POST")
	r.HandleFunc("/pages/{id}", GetPage)
	http.Handle("/", r)
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, world.")
}
