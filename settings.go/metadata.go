package utils

import (
	m "fluorescences/models"
)

// GetMetadata will return a metadata struct from the settings bucket
func GetMetadata() (meta m.Metadata, err error) {

	meta = m.Metadata{
		Title: "Fluorescences",
		Desc:  "A comic blog",
	}

	return
}
