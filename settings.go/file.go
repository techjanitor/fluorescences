package utils

import (
	"io"
	"mime/multipart"
	"os"

	m "fluorescences/models"
)

// SaveFile saves an uploaded file to disk
func SaveFile(file multipart.File, filename string) (filetype m.FileType, err error) {

	var dst *os.File

	dst, err = os.Create("images/" + filename)
	if err != nil {
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		return
	}

	filetype = m.FileType{
		Filename: filename,
	}

	return

}
