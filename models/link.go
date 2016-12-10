package models

import "net/url"

// LinkType holds a key for an external link
type LinkType struct {
	ID      int
	Title   string `storm:"unique"`
	Address string
	URL     *url.URL
}

// Links is a slice of LinkType
type Links []LinkType

func (f Links) Len() int {
	return len(f)
}

func (f Links) Less(i, j int) bool {
	return f[i].ID < f[j].ID
}

func (f Links) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}
