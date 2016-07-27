package models

import (
	"html/template"
	"time"
)

// BlogType holds a blog post
type BlogType struct {
	ID         int
	User       string
	Title      string `storm:"unique"`
	Content    string
	ContentOut template.HTML
	HumanTime  string
	StoredTime time.Time
}
