package utils

import (
	m "fluorescences/models"
	"fmt"
	"time"
)

const (
	defaultCommission = `# Commission Info

## This supports markdown!

1. Rules
1. Something
2. Something Else

- Need money
- Send money

### A table

| Price        | Type           | Shading  |
| ------------- |:-------------:| -----:|
| $100      | sketch | nope |`
)

// GetMetadata will return a metadata struct from the settings bucket
func GetMetadata() (meta m.Metadata, err error) {

	err = Storm.Get("data", "metadata", &meta)
	if err != nil {
		return
	}

	// get all the links
	err = Storm.All(&meta.Links)
	if err != nil {
		return
	}

	return
}

// InitData will set the initial data
func InitData(name string) (err error) {

	// init the store
	Initialize(name)

	tx, err := Storm.Begin(true)
	if err != nil {
		return
	}
	defer tx.Rollback()

	com := m.CommissionType{
		Open:        false,
		Content:     defaultCommission,
		UpdatedTime: time.Now(),
	}

	err = tx.Set("data", "commission", &com)
	if err != nil {
		return
	}

	meta := m.Metadata{
		Title: "Fluorescences",
		Desc:  "A comic blog",
	}

	err = tx.Set("data", "metadata", &meta)
	if err != nil {
		return
	}

	cat := m.CategoryType{
		Title: "Default",
		Desc:  "The default category",
	}

	// save gallery category
	err = tx.Save(&cat)
	if err != nil {
		return
	}

	blog := m.BlogType{
		User:       MustGetUsername(),
		StoredTime: time.Now(),
		Title:      "First post",
		Content:    "Welcome to your new art blog!",
	}

	err = tx.Save(&blog)
	if err != nil {
		return
	}

	link := m.LinkType{
		Title:   "Tumblr",
		Address: "http://tumblr.com",
	}

	// save gallery
	err = tx.Save(&link)
	if err != nil {
		return
	}

	// commit
	tx.Commit()

	fmt.Println("Data Generated")

	return
}
