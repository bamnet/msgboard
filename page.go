package msgboard

import (
	"appengine"
	"appengine/datastore"
	"fmt"
	"net/http"
	"time"

	"encoding/json"
	"github.com/gorilla/mux"
)

type Page struct {
	Title       string
	Content     string `datastore:",noindex"`
	LastUpdated time.Time
	ID string `datastore:"-"`
}

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

	b, _ := json.Marshal(pages)
	fmt.Fprint(w, string(b))
}

func CreatePage(w http.ResponseWriter, r *http.Request){
	c := appengine.NewContext(r)
	decoder := json.NewDecoder(r.Body)
	var p Page   
	if err := decoder.Decode(&p); err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
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

	b, _ := json.Marshal(page)
	fmt.Fprint(w, string(b))
}
