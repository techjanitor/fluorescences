package models

import "time"

// KeyType holds a key for a private gallery
type KeyType struct {
	ID      int
	Key     string
	Created time.Time
	Expires bool
}

// Keys is a slice of KeyType
type Keys []KeyType

func (f Keys) Len() int {
	return len(f)
}

func (f Keys) Less(i, j int) bool {
	return f[i].ID < f[j].ID
}

func (f Keys) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}
