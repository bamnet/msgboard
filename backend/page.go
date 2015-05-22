package msgboard

import (
	"appengine"
	"appengine/datastore"
	"appengine/memcache"
	"errors"
	"io"
	"time"

	"encoding/json"
	"github.com/russross/blackfriday"
)

// Page models a slide to be shown on the message board.
type Page struct {
	Title       string    `json:"title"`                        // Title of the Page
	Content     string    `json:"content" datastore:",noindex"` // Markdown text content.
	LastUpdated time.Time `json:"last_updated"`                 // Timestamp the content was last updated.
	ID          string    `json:"id" datastore:"-"`             // Encoded datastore key acting as an ID.
	Rendered    string    `json:"rendered" datastore:"-"`       // HTML for the page based rendered from Content.
}

// Error codes returned for invalid pages.
var (
	ErrMissingTitle = errors.New("page missing title") // Pages always need a title.
)

// ListPages returns all the pages.
// The view arg can be used to restrict the output level:
//   (blank) - Full results
//   ids - Results with only the ID attribute.
func ListPages(ctx appengine.Context, view string) ([]Page, error) {
	q := datastore.NewQuery("Page")
	q = q.Order("Title")
	var pages []Page

	if view == "ids" {
		q = q.KeysOnly()
	}
	keys, err := q.GetAll(ctx, &pages)
	if err != nil {
		return nil, err
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

	return pages, nil
}

// CreatePage uses JSON data to create and save a new Page.
func CreatePage(ctx appengine.Context, jsonBody io.Reader) (*Page, error) {
	var p Page
	if err := jsonToPage(jsonBody, &p); err != nil {
		return nil, err
	}
	p.LastUpdated = time.Now()

	key, err := datastore.Put(ctx, datastore.NewIncompleteKey(ctx, "Page", nil), &p)
	if err != nil {
		return nil, err
	}
	p.ID = key.Encode()
	return &p, nil
}

// GetPage uses a Page ID from the URL to return a Page.
func GetPage(ctx appengine.Context, ID string) (*Page, error) {
	page, _ := getPageMemcache(ctx, ID)
	if page == nil {
		ctx.Debugf("memcache miss for id: %s", ID)
		var err error
		page, err = getPageDatastore(ctx, ID)
		if err != nil {
			return nil, err
		}
		defer setPageMemcache(ctx, page)
	}

	page.Rendered = string(blackfriday.MarkdownCommon([]byte(page.Content)))
	return page, nil
}

func getPageDatastore(ctx appengine.Context, ID string) (*Page, error) {
	k, err := datastore.DecodeKey(ID)
	if err != nil {
		return nil, err
	}

	var page Page
	if err := datastore.Get(ctx, k, &page); err != nil {
		return nil, err
	}
	page.ID = k.Encode()
	return &page, nil
}

func getPageMemcache(ctx appengine.Context, ID string) (*Page, error) {
	var page Page
	if _, err := memcache.JSON.Get(ctx, ID, &page); err != nil {
		return nil, err
	}
	return &page, nil
}

func setPageMemcache(ctx appengine.Context, page *Page) {
	// Don't store the rendered HTML in memcache
	page.Rendered = ""

	i := &memcache.Item{
		Key:    page.ID,
		Object: page,
	}

	if err := memcache.JSON.Set(ctx, i); err != nil {
		ctx.Errorf("error setting memcache: %s", err)
	}
}

// UpdatePage uses a Page ID from the URL and JSON data to update a Page.
func UpdatePage(ctx appengine.Context, ID string, jsonBody io.Reader) (*Page, error) {
	key, err := datastore.DecodeKey(ID)
	if err != nil {
		return nil, err
	}

	var page Page
	if err := datastore.Get(ctx, key, &page); err != nil {
		return nil, err
	}
	if err := jsonToPage(jsonBody, &page); err != nil {
		return nil, err
	}
	page.ID = ID
	page.LastUpdated = time.Now()

	if _, err := datastore.Put(ctx, key, &page); err != nil {
		return nil, err
	}
	defer setPageMemcache(ctx, &page)
	return &page, nil
}

// DeletePage uses a Page ID to remove a page from the datastore.
func DeletePage(ctx appengine.Context, ID string) error {
	key, err := datastore.DecodeKey(ID)
	if err != nil {
		return err
	}
	return datastore.Delete(ctx, key)
}

func jsonToPage(data io.Reader, page *Page) error {
	decoder := json.NewDecoder(data)
	if err := decoder.Decode(&page); err != nil {
		return err
	}
	if page.Title == "" {
		return ErrMissingTitle
	}
	return nil
}
