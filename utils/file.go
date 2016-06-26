package utils

import (
	"io"
	"mime/multipart"
	"os"
)

// FileType holds an image file
type FileType struct {
	ID       int
	Filename string
}

// SaveFile saves an uploaded file to disk
func SaveFile(file multipart.File, filename string) (filetype FileType, err error) {

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

	filetype = FileType{
		Filename: filename,
	}

	return

}
