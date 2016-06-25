package utils

import (
	"github.com/boltdb/bolt"
)

var (
	// Bolt holds the database handle
	Bolt *bolt.DB
)

func init() {
	var err error

	Bolt, err = bolt.Open("bolt.db", 0600, nil)
	if err != nil {
		panic(err)
	}

}
