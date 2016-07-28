package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	m "fluorescences/models"
	u "fluorescences/utils"
)

type newForm struct {
	Title string `form:"title" binding:"required"`
	Desc  string `form:"desc" binding:"required"`
}

// PostController posts new galleries
func PostController(c *gin.Context) {
	var err error
	var nf newForm

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
		c.Error(err).SetMeta("gallery.PostController.FormFile")
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

	// start transaction
	tx, err := u.Storm.Begin(true)
	if err != nil {
		c.Error(err).SetMeta("gallery.PostController.Begin")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}
	defer tx.Rollback()

	// save gallery
	err = tx.Save(&gallery)
	if err != nil {
		c.Error(err).SetMeta("gallery.PostController.Save")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	notification := m.BlogType{
		User:          "test",
		Notificiation: true,
		StoredTime:    time.Now(),
		Title:         "New Gallery",
		Content:       fmt.Sprintf("<a href=\"/comic/%d/1\">%s</a>", gallery.ID, gallery.Title),
	}

	// save blog notification
	err = tx.Save(&notification)
	if err != nil {
		c.Error(err).SetMeta("gallery.PostController.Save")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	// commit
	tx.Commit()

	c.Redirect(http.StatusFound, "/comics/1")

	return

}
