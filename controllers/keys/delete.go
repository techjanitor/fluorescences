package controllers

import (
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"

	m "fluorescences/models"
	u "fluorescences/utils"
)

type deleteForm struct {
	Gallery int    `form:"gallery" binding:"required"`
	Key     string `form:"key" binding:"required"`
}

// DeleteController deletes keys from galleries
func DeleteController(c *gin.Context) {
	var err error
	var df deleteForm

	err = c.Bind(&df)
	if err != nil {
		c.Error(err).SetMeta("keys.DeleteController.Bind")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	// start transaction
	tx, err := u.Storm.Begin(true)
	if err != nil {
		c.Error(err).SetMeta("keys.DeleteController.Begin")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}
	defer tx.Rollback()

	var gallery m.GalleryType

	err = tx.One("ID", df.Gallery, &gallery)
	if err != nil {
		c.Error(err).SetMeta("keys.DeleteController.One")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	sort.Sort(gallery.Keys)

	// remove the file from the slice
	for i := len(gallery.Keys) - 1; i >= 0; i-- {
		key := gallery.Keys[i]

		if key.Key == df.Key {
			gallery.Keys = append(gallery.Keys[:i], gallery.Keys[i+1:]...)
		}
	}

	err = tx.Save(&gallery)
	if err != nil {
		c.Error(err).SetMeta("keys.DeleteController.Save")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	// commit
	tx.Commit()

	c.Redirect(http.StatusFound, c.Request.Referer())

	return

}
