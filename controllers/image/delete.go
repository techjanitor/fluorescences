package controllers

import (
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"

	m "fluorescences/models"
	u "fluorescences/utils"
)

type deleteForm struct {
	Gallery int `form:"gallery" binding:"required"`
	Image   int `form:"image" binding:"required"`
}

// DeleteController deletes images from galleries
func DeleteController(c *gin.Context) {
	var err error
	var df deleteForm

	err = c.Bind(&df)
	if err != nil {
		c.Error(err).SetMeta("image.DeleteController.Bind")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	// start transaction
	tx, err := u.Storm.Begin(true)
	if err != nil {
		c.Error(err).SetMeta("image.DeleteController.Begin")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}
	defer tx.Rollback()

	var gallery m.GalleryType

	err = tx.One("ID", df.Gallery, &gallery)
	if err != nil {
		c.Error(err).SetMeta("image.DeleteController.One")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	sort.Sort(gallery.Files)

	// remove the file from the slice
	for i := len(gallery.Files) - 1; i >= 0; i-- {
		file := gallery.Files[i]

		if file.ID == df.Image {
			gallery.Files = append(gallery.Files[:i], gallery.Files[i+1:]...)
		}
	}

	err = tx.Save(&gallery)
	if err != nil {
		c.Error(err).SetMeta("image.DeleteController.Save")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	// commit
	tx.Commit()

	c.Redirect(http.StatusFound, c.Request.Referer())

	return

}
