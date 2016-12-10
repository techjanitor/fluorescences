package models

import (
	"html/template"
	"time"
)

// CategoryType holds category info
type CategoryType struct {
	ID        int
	Galleries int
	Cover     string
	Title     string `storm:"unique"`
	Desc      string
	DescOut   template.HTML
}

// Categories is a slice of Categorys
type Categories []CategoryType

// GalleryType holds a gallery
type GalleryType struct {
	ID          int
	Images      int
	Category    int    `storm:"index"`
	Title       string `storm:"unique"`
	Cover       string
	Desc        string
	DescOut     template.HTML
	Private     bool `storm:"index"`
	HumanTime   string
	StoredTime  time.Time `storm:"index"`
	UpdatedTime time.Time `storm:"index"`
	Files       Files
	Keys        Keys
}

// Galleries is a slice of GalleryTypes
type Galleries []GalleryType
