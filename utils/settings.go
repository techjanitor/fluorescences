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
	return
}

// InitData will set the initial data
func InitData() (err error) {

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

	// commit
	tx.Commit()

	fmt.Println("Data Generated")

	return
}
