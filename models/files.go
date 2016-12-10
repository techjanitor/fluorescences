package models

import "html/template"

// FileType holds an image file
type FileType struct {
	ID       int
	Filename string
	Desc     string
	DescOut  template.HTML
}

// Files is a slice of FileTypes
type Files []FileType

func (f Files) Len() int {
	return len(f)
}

func (f Files) Less(i, j int) bool {
	return f[i].ID < f[j].ID
}

func (f Files) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}
