package controllers

import (
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"

	m "fluorescences/models"
	u "fluorescences/utils"
)

type newForm struct {
	ID int `form:"id" binding:"required"`
}

// NewController add an image to a gallery
func NewController(c *gin.Context) {
	var err error
	var inf newForm

	err = c.Bind(&inf)
	if err != nil {
		c.Error(err).SetMeta("image.NewController.Bind")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	// Check if theres a file
	upload, fileheader, err := c.Request.FormFile("file")
	if err != nil {
		c.Error(err).SetMeta("image.NewController.FormFile")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	// save the file to disk
	file, err := u.SaveFile(upload, fileheader.Filename)
	if err != nil {
		c.Error(err).SetMeta("image.NewController.SaveFile")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	// add the image to the gallery file slice
	err = AddImage(inf.ID, file)
	if err != nil {
		c.Error(err).SetMeta("image.NewController.AddImage")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	c.Redirect(http.StatusFound, c.Request.Referer())

	return

}

// AddImage will add a blog post
func AddImage(gid int, file m.FileType) (err error) {

	// start transaction
	tx, err := u.Storm.Begin(true)
	if err != nil {
		return
	}
	defer tx.Rollback()

	var gallery m.GalleryType

	err = tx.One("ID", gid, &gallery)
	if err != nil {
		return
	}

	sort.Sort(gallery.Files)

	// set the id for the file
	if len(gallery.Files) == 0 {
		// if this is the first file
		file.ID = 1
	} else {
		// use the next sequential number
		lid := gallery.Files[len(gallery.Files)-1]
		file.ID = lid.ID + 1
	}

	gallery.Files = append(gallery.Files, file)

	err = tx.Save(&gallery)
	if err != nil {
		return
	}

	// commit
	tx.Commit()

	return

}
