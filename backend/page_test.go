package msgboard

import (
	"reflect"
	"sort"
	"testing"
	"time"

	"appengine"
	"appengine/aetest"
	"appengine/datastore"
)

func TestGetMissingPage(t *testing.T) {
	ctx, err := aetest.NewContext(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer ctx.Close()

	ID := datastore.NewKey(ctx, "Page", "missing", 0, nil).Encode()
	if _, err := GetPage(ctx, ID); err != datastore.ErrNoSuchEntity {
		t.Errorf("want: NoSuchEntity, got: %v", err)
	}
}

func TestGetPage(t *testing.T) {
	ctx, err := aetest.NewContext(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer ctx.Close()

	p := Page{Title: "FooBar", Content: "Body here"}
	key, err := datastore.Put(ctx, datastore.NewIncompleteKey(ctx, "Page", nil), &p)
	if err != nil {
		t.Fatal(err)
	}
	ID := key.Encode()

	page, err := GetPage(ctx, ID)
	if err != nil {
		t.Fatal(err)
	}
	if page.Title != p.Title && page.Content != p.Content {
		t.Error("Unexpected page")
	}
}

// ByTitle implements sort.Interface for []Page based on
// the Title field.
type ByTitle []Page

func (a ByTitle) Len() int           { return len(a) }
func (a ByTitle) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByTitle) Less(i, j int) bool { return a[i].Title < a[j].Title }

func TestListPages(t *testing.T) {
	ctx, err := aetest.NewContext(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer ctx.Close()

	p, _, err := setupPages(ctx)

	r, err := ListPages(ctx, "")
	if err != nil {
		t.Fatal(err)
	}
	if got, want := len(r), len(p); got != want {
		t.Errorf("got len: %d want len: %d", got, want)
	}

	sort.Sort(ByTitle(p))
	sort.Sort(ByTitle(r))
	for i := range p {
		if gotT, gotC, wantT, wantC := p[i].Title, p[i].Content, r[i].Title, r[i].Content; gotT != wantT || gotC != wantC {
			t.Errorf("got title %s content %s, want title %s, content %s", gotT, gotC, wantT, wantC)
		}
	}
}

func TestListPagesViewIDs(t *testing.T) {
	ctx, err := aetest.NewContext(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer ctx.Close()

	_, k, err := setupPages(ctx)

	r, err := ListPages(ctx, "ids")
	if err != nil {
		t.Fatal(err)
	}
	if got, want := len(r), len(k); got != want {
		t.Errorf("got len: %d want len: %d", got, want)
	}

	var gotIDs []string
	var wantIDs []string
	for i := range k {
		gotIDs = append(gotIDs, r[i].ID)
		wantIDs = append(wantIDs, k[i].Encode())
	}
	sort.Strings(gotIDs)
	sort.Strings(wantIDs)
	if !reflect.DeepEqual(gotIDs, wantIDs) {
		t.Errorf("got %s, want: %s", gotIDs, wantIDs)
	}
}

func setupPages(ctx appengine.Context) ([]Page, []*datastore.Key, error) {
	p := []Page{
		{Title: "1", Content: "Body here"},
		{Title: "2", Content: ""},
	}
	k, err := datastore.PutMulti(ctx, []*datastore.Key{
		datastore.NewIncompleteKey(ctx, "Page", nil),
		datastore.NewIncompleteKey(ctx, "Page", nil),
	}, p)
	if err != nil {
		return nil, nil, err
	}
	time.Sleep(1 * time.Second)

	return p, k, nil
}
