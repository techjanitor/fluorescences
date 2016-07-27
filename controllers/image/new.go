package controllers

import (
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"

	m "fluorescences/models"
	u "fluorescences/utils"
)

// NewForm is the input from the new image form
type NewForm struct {
	ID int `form:"id" binding:"required"`
}

// NewController add an image to a gallery
func NewController(c *gin.Context) {
	var err error
	var inf NewForm

	err = c.Bind(&inf)
	if err != nil {
		c.Error(err).SetMeta("image.NewController.Bind")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	// Check if theres a file
	upload, fileheader, err := c.Request.FormFile("file")
	if err != nil {
		c.Error(err).SetMeta("image.NewController")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	file, err := u.SaveFile(upload, fileheader.Filename)
	if err != nil {
		c.Error(err).SetMeta("image.NewController.SaveFile")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

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

	var gallery m.GalleryType

	u.Storm.One("ID", gid, &gallery)

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

	u.Storm.Save(&gallery)

	return

}
