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

// Page models a slide to be shown on the message board.
type Page struct {
	Title       string    // Title of the Page
	Content     string    `datastore:",noindex"` // Markdown text content.
	LastUpdated time.Time // Timestamp the content was last updated.
	ID          string    `datastore:"-"` // Encoded datastore key acting as an ID.
	Rendered    string    `datastore:"-"` // HTML for the page based rendered from Content.
}

// Error codes returned for invalid pages.
var (
	ErrMissingTitle = errors.New("page missing title") // Pages always need a title.
)

// ListPages returns JSON of all the pages.
// The view GET parameter can be used to restrict the output level:
//   (blank / default) - Full results
//   ids - Results with only the ID attribute.
func ListPages(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	q := datastore.NewQuery("Page")
	var pages []Page

	if r.FormValue("view") == "ids" {
		q = q.KeysOnly()
	}
	keys, err := q.GetAll(c, &pages)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(keys) > len(pages) {
		// Handle instances where not enough pages were created for each key.
		// This behavior happens when GetAll doesn't create pages like a KeyOnly query.
		diff := len(keys) - len(pages)
		pages = append(pages, make([]Page, diff)...)
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

// CreatePage uses JSON POST data to create and save a new Page.
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

// GetPage uses a Page ID from the URL to return JSON of the Page.
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

// UpdatePage uses a Page ID from the URL and JSON POST data to Update a Page.
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
