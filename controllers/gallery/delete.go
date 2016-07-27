package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	m "fluorescences/models"
	u "fluorescences/utils"
)

// DeleteForm takes the input from the delete form
type DeleteForm struct {
	Gallery int `form:"gallery" binding:"required"`
}

// DeleteController deletes galleries
func DeleteController(c *gin.Context) {
	var err error
	var df DeleteForm

	err = c.Bind(&df)
	if err != nil {
		c.Error(err).SetMeta("gallery.DeleteController.Bind")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	var gallery m.GalleryType

	// get the gallery from bolt
	err = u.Storm.One("ID", df.Gallery, &gallery)
	if err != nil {
		c.Error(err).SetMeta("gallery.DeleteController.One")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	// delete it
	err = u.Storm.Remove(&gallery)
	if err != nil {
		c.Error(err).SetMeta("gallery.DeleteController.Remove")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	c.Redirect(http.StatusFound, c.Request.Referer())

	return

}
