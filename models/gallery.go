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
	Keys        Keys
}

// GalleryCategory holds sorted galleries
type GalleryCategory struct {
	ID        int
	Title     string `storm:"unique"`
	Desc      string
	Galleries Galleries
}

// Galleries is a slice of GalleryTypes
type Galleries []*GalleryType
