// Package msgboard provides a message board application for users to
// create and edit HTML pages shown on a screen.
package msgboard

import (
	"appengine"
	"appengine/datastore"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func init() {
	r := mux.NewRouter()
	r.HandleFunc("/api/pages", listPagesHandler).Methods("GET")
	r.HandleFunc("/api/pages", createPageHandler).Methods("POST")
	r.HandleFunc("/api/pages/{id}", getPageHandler).Methods("GET")
	r.HandleFunc("/api/pages/{id}", updatePageHandler).Methods("PATCH")
	r.HandleFunc("/api/pages/{id}", deletePageHandler).Methods("DELETE")
	http.Handle("/api/", r)

	http.HandleFunc("/", homeRedirectHandler)
}

func writeJSON(w http.ResponseWriter, v interface{}) {
	b, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
	return
}

func writeError(w http.ResponseWriter, err error) {
	ec := http.StatusInternalServerError
	switch err {
	case datastore.ErrNoSuchEntity:
		ec = http.StatusNotFound
	case ErrMissingTitle:
		ec = http.StatusBadRequest
	}
	http.Error(w, err.Error(), ec)
	return
}

func homeRedirectHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/admin/", http.StatusMovedPermanently)
	return
}

func listPagesHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	view := r.FormValue("view")

	pages, err := ListPages(ctx, view)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, pages)
	return
}

func createPageHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	page, err := CreatePage(ctx, r.Body)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, page)
	return
}

func getPageHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	vars := mux.Vars(r)
	ID := vars["id"]

	page, err := GetPage(ctx, ID)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, page)
	return
}

func updatePageHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	vars := mux.Vars(r)
	ID := vars["id"]

	page, err := UpdatePage(ctx, ID, r.Body)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, page)
	return
}

func deletePageHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	vars := mux.Vars(r)
	ID := vars["id"]

	if err := DeletePage(ctx, ID); err != nil {
		writeError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	return
}
