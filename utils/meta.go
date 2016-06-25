package utils

import (
	"encoding/json"

	"github.com/boltdb/bolt"
)

var (
	// SettingsDB holds the name of the settings bucket
	SettingsDB = "settings"
)

// Metadata holds the meta data for pages
type Metadata struct {
	Title string `json:"page_title"`
	Desc  string `json:"page_description"`
}

// GetMetadata will return a metadata struct from the settings bucket
func GetMetadata() (meta Metadata, err error) {
	err = Bolt.View(func(tx *bolt.Tx) (err error) {
		// the settings bucket
		s := tx.Bucket([]byte(SettingsDB))

		// stats for key count
		stats := s.Stats()

		if stats.KeyN < 1 {
			meta = Metadata{}
			return
		}

		cs := s.Cursor()

		_, settings := cs.Seek(Itob(1))
		err = json.Unmarshal(settings, &meta)
		if err != nil {
			return
		}

		return
	})
	if err != nil {
		return
	}

	return
}

// InitMetadata will create a new metadata bucket with the default settings
func InitMetadata() (err error) {
	err = Bolt.Update(func(tx *bolt.Tx) (err error) {
		bucket, err := tx.CreateBucketIfNotExists([]byte(SettingsDB))
		if err != nil {
			return
		}

		// default info
		metadata := Metadata{
			Title: "Fluorescences",
			Desc:  "A comic blog",
		}

		// encode our roomconfig
		encoded, err := json.Marshal(metadata)
		if err != nil {
			return
		}

		// put the blog post
		err = bucket.Put(Itob(1), encoded)
		if err != nil {
			return
		}

		return
	})
	if err != nil {
		return
	}

	return
}
