package models

import (
	u "fluorescences/utils"
)

// Files is a slice of FileTypes
type Files []u.FileType

func (f Files) Len() int {
	return len(f)
}

func (f Files) Less(i, j int) bool {
	return f[i].ID < f[j].ID
}

func (f Files) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}
