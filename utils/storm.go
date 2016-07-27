package utils

import "github.com/asdine/storm"

var (
	// Storm holds the database handle
	Storm *storm.DB
)

func init() {
	var err error

	// open the database
	Storm, err = storm.Open("data.db", storm.AutoIncrement())
	if err != nil {
		panic(err)
	}

}
