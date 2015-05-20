// Package msgboard provides a message board application for users to
// create and edit HTML pages shown on a screen.
package msgboard

import (
	"appengine"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func init() {
	r := mux.NewRouter()
	r.HandleFunc("/api/pages", listPagesHandler).Methods("GET")
	r.HandleFunc("/api/pages", createPageHandler).Methods("POST")
	r.HandleFunc("/api/pages/{id}", getPageHandler).Methods("GET")
	r.HandleFunc("/api/pages/{id}", updatePageHandler).Methods("PATCH")
	http.Handle("/api/", r)
}

func listPagesHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	view := r.FormValue("view")

	pages, err := ListPages(ctx, view)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	b, _ := json.Marshal(pages)
	fmt.Fprint(w, string(b))
	return
}

func createPageHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	page, err := CreatePage(ctx, r.Body)
	if err != nil {
		ec := http.StatusInternalServerError
		if err == ErrMissingTitle {
			ec = http.StatusBadRequest
		}
		http.Error(w, err.Error(), ec)
		return
	}

	b, _ := json.Marshal(page)
	fmt.Fprint(w, string(b))
	return
}

func getPageHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	vars := mux.Vars(r)
	ID := vars["id"]

	page, err := GetPage(ctx, ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	b, _ := json.Marshal(page)
	fmt.Fprint(w, string(b))
	return
}

func updatePageHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	vars := mux.Vars(r)
	ID := vars["id"]

	page, err := UpdatePage(ctx, ID, r.Body)
	if err != nil {
		ec := http.StatusInternalServerError
		if err == ErrMissingTitle {
			ec = http.StatusBadRequest
		}
		http.Error(w, err.Error(), ec)
		return
	}

	b, _ := json.Marshal(page)
	fmt.Fprint(w, string(b))
	return
}
