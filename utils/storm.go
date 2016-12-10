package utils

import (
	"fmt"

	"github.com/asdine/storm"
)

var (
	// Storm holds the database handle
	Storm *storm.DB
)

// Initialize will open the store
func Initialize(name string) {
	var err error

	// open the database
	Storm, err = storm.Open(fmt.Sprintf("%s.db", name), storm.AutoIncrement())
	if err != nil {
		panic(err)
	}

}
