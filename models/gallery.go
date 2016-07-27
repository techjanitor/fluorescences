package models

import (
	"time"
)

// GalleryType holds a gallery
type GalleryType struct {
	ID         int
	User       string
	Title      string
	Cover      string
	Desc       string
	Private    bool
	HumanTime  string
	OpenTime   time.Time
	StoredTime time.Time
	Files      Files
	Keys       []KeyType
}
