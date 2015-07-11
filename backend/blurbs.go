package msgboard

import (
	"appengine"
	"appengine/datastore"
	"appengine/memcache"
	"io"
	"time"

	"encoding/json"
)

// Blurbs store a map of fields and strings to display in those fields.
type Blurbs struct {
	Content        map[string]string `json:"content" datastore:"-"`
	LastUpdated    time.Time         `json:"last_updated"` // Timestamp the content was last updated.
	EncodedContent []byte            `json:"-"`
}

// Load unpacks the blurbs by converting the encoded content into the contents map.
func (b *Blurbs) Load(c <-chan datastore.Property) error {
	if err := datastore.LoadStruct(b, c); err != nil {
		return err
	}

	if err := json.Unmarshal(b.EncodedContent, &b.Content); err != nil {
		return err
	}
	return nil
}

// Save packs the blurb for storage by encoding the contents map into something storable.
func (b *Blurbs) Save(c chan<- datastore.Property) error {
	var err error
	if b.EncodedContent, err = json.Marshal(b.Content); err != nil {
		return err
	}
	if err := datastore.SaveStruct(b, c); err != nil {
		return err
	}
	return nil
}

// GetBlurbs returns the current Blurbs. If none exist, creates them.
func GetBlurbs(ctx appengine.Context) (*Blurbs, error) {
	blurbs, _ := getBlurbsMemcache(ctx)
	if blurbs == nil {
		ctx.Debugf("memcache miss for blurbs")
		var err error
		blurbs, err = getBlurbsDatastore(ctx)
		if err != nil {
			if err == datastore.ErrNoSuchEntity {
				return initBlurbs(ctx)
			}
			return nil, err
		}
		defer setBlurbsMemcache(ctx, blurbs)
	}
	return blurbs, nil
}

// UpdateBlurbs uses JSON data to update the Blurbs.
func UpdateBlurbs(ctx appengine.Context, jsonBody io.Reader) (*Blurbs, error) {
	key := blurbsKey(ctx)

	var blurbs Blurbs
	decoder := json.NewDecoder(jsonBody)
	if err := decoder.Decode(&blurbs); err != nil {
		return nil, err
	}
	blurbs.LastUpdated = time.Now()

	if _, err := datastore.Put(ctx, key, &blurbs); err != nil {
		return nil, err
	}
	defer setBlurbsMemcache(ctx, &blurbs)
	return &blurbs, nil
}

// initBlurbs creates initial (empty) Blurbs.
func initBlurbs(ctx appengine.Context) (*Blurbs, error) {
	var b Blurbs
	b.LastUpdated = time.Now()

	_, err := datastore.Put(ctx, blurbsKey(ctx), &b)
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func getBlurbsDatastore(ctx appengine.Context) (*Blurbs, error) {
	var blurbs Blurbs
	if err := datastore.Get(ctx, blurbsKey(ctx), &blurbs); err != nil {
		return nil, err
	}
	return &blurbs, nil
}

func getBlurbsMemcache(ctx appengine.Context) (*Blurbs, error) {
	var blurbs Blurbs
	if _, err := memcache.JSON.Get(ctx, blurbsKey(ctx).String(), &blurbs); err != nil {
		return nil, err
	}
	return &blurbs, nil
}

func setBlurbsMemcache(ctx appengine.Context, blurbs *Blurbs) {
	i := &memcache.Item{
		Key:    blurbsKey(ctx).String(),
		Object: blurbs,
	}

	if err := memcache.JSON.Set(ctx, i); err != nil {
		ctx.Errorf("error setting memcache: %s", err)
	}
}

// blurbsKey returns the datastore key used for the Blurbs.
// Since there is only 1 set of blurbs they key is static.
func blurbsKey(ctx appengine.Context) *datastore.Key {
	return datastore.NewKey(ctx, "Blurbs", "blurbs", 0, nil)
}
