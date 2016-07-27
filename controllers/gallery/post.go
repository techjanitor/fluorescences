package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	m "fluorescences/models"
	u "fluorescences/utils"
)

// NewForm is the input from the gallery form
type NewForm struct {
	Title string `form:"title" binding:"required"`
	Desc  string `form:"desc" binding:"required"`
}

// PostController posts new blogs
func PostController(c *gin.Context) {
	var err error
	var nf NewForm

	err = c.Bind(&nf)
	if err != nil {
		c.Error(err).SetMeta("gallery.PostController.Bind")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	// slice to hold our file information
	var files m.Files

	// Check if theres a file
	upload, fileheader, err := c.Request.FormFile("file")
	if err != nil {
		c.Error(err).SetMeta("gallery.PostController")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	file, err := u.SaveFile(upload, fileheader.Filename)
	if err != nil {
		c.Error(err).SetMeta("gallery.PostController.SaveFile")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	// this is always the first image in the gallery
	file.ID = 1

	files = append(files, file)

	gallery := m.GalleryType{
		User:       "test",
		StoredTime: time.Now(),
		Title:      nf.Title,
		Desc:       nf.Desc,
		Files:      files,
	}

	err = u.Storm.Save(&gallery)
	if err != nil {
		c.Error(err).SetMeta("gallery.PostController.Save")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	c.Redirect(http.StatusFound, "/comics/1")

	return

}
