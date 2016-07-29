package models

import (
	"time"
)

// GalleryType holds a gallery
type GalleryType struct {
	ID          int
	Title       string `storm:"unique"`
	Cover       string
	Desc        string
	Private     bool `storm:"index"`
	HumanTime   string
	StoredTime  time.Time `storm:"index"`
	UpdatedTime time.Time `storm:"index"`
	Files       Files
	Keys        []KeyType
}
