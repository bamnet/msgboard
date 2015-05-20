package msgboard

import (
	"testing"

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
