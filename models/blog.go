package models

import (
	"html/template"
	"time"
)

// BlogType holds a blog post
type BlogType struct {
	ID            int
	Notificiation bool
	User          string
	Title         string
	Content       string
	ContentOut    template.HTML
	HumanTime     string
	StoredTime    time.Time `storm:"index"`
}

// Blogs is a slice of BlogTypes
type Blogs []BlogType
