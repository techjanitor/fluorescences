package models

import "time"

// KeyType holds a key for a private gallery
type KeyType struct {
	ID      int
	Key     string
	Created time.Time
}
