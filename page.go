package msgboard

import (
	"appengine"
	"appengine/datastore"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/russross/blackfriday"
)

type Page struct {
	Title       string
	Content     string `datastore:",noindex"`
	LastUpdated time.Time
	ID          string `datastore:"-"`
	Rendered    string `datastore:"-"`
}

var (
	ErrMissingTitle = errors.New("page missing title")
)

func ListPages(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	q := datastore.NewQuery("Page")
	var pages []Page

	keys, err := q.GetAll(c, &pages)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for i, k := range keys {
		pages[i].ID = k.Encode()
	}

	if len(pages) == 0 {
		fmt.Fprint(w, "[]")
		return
	}
	b, _ := json.Marshal(pages)
	fmt.Fprint(w, string(b))
	return
}

func CreatePage(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	var p Page
	if err, ec := jsonToPage(r.Body, &p); err != nil {
		http.Error(w, err.Error(), ec)
		return
	}
	p.LastUpdated = time.Now()

	key, err := datastore.Put(c, datastore.NewIncompleteKey(c, "Page", nil), &p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	p.ID = key.Encode()
	b, _ := json.Marshal(p)
	fmt.Fprint(w, string(b))
}

func GetPage(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	vars := mux.Vars(r)
	ID := vars["id"]

	k, err := datastore.DecodeKey(ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var page Page

	if err := datastore.Get(c, k, &page); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	page.ID = k.Encode()
	page.Rendered = string(blackfriday.MarkdownCommon([]byte(page.Content)))

	b, _ := json.Marshal(page)
	fmt.Fprint(w, string(b))
}

func UpdatePage(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	vars := mux.Vars(r)
	ID := vars["id"]

	key, err := datastore.DecodeKey(ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var page Page

	if err := datastore.Get(c, key, &page); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err, ec := jsonToPage(r.Body, &page); err != nil {
		http.Error(w, err.Error(), ec)
		return
	}
	page.ID = ID
	page.LastUpdated = time.Now()

	if _, err := datastore.Put(c, key, &page); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	b, _ := json.Marshal(page)
	fmt.Fprint(w, string(b))
}

func jsonToPage(data io.Reader, page *Page) (error, int) {
	decoder := json.NewDecoder(data)
	if err := decoder.Decode(&page); err != nil {
		return err, http.StatusInternalServerError
	}
	if page.Title == "" {
		return ErrMissingTitle, http.StatusBadRequest
	}
	return nil, http.StatusOK
}
