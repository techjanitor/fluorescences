package utils

import (
	"github.com/boltdb/bolt"
)

var (
	// Bolt holds the database handle
	Bolt *bolt.DB
	// SettingsDB holds the name of the settings bucket
	SettingsDB = "settings"
	//GalleryDB is the bucket for comic galleries
	GalleryDB = "galleries"
	// BlogDB is the bucket name for blog posts
	BlogDB = "blogs"
)

func init() {
	var err error

	// open the database
	Bolt, err = bolt.Open("bolt.db", 0600, nil)
	if err != nil {
		panic(err)
	}

	// create our buckets if they dont exist yet
	err = Bolt.Update(func(tx *bolt.Tx) (err error) {
		_, err = tx.CreateBucketIfNotExists([]byte(GalleryDB))
		if err != nil {
			return
		}

		_, err = tx.CreateBucketIfNotExists([]byte(BlogDB))
		if err != nil {
			return
		}

		_, err = tx.CreateBucketIfNotExists([]byte(SettingsDB))
		if err != nil {
			return
		}

		return
	})
	if err != nil {
		panic("could not init buckets")
	}

}
