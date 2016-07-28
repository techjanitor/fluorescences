package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	m "fluorescences/models"
	u "fluorescences/utils"
)

type updateForm struct {
	ID    int    `form:"id" binding:"required"`
	Title string `form:"title" binding:"required"`
	Desc  string `form:"desc" binding:"required"`
}

// UpdateController updates gallery information
func UpdateController(c *gin.Context) {
	var err error
	var uf updateForm

	err = c.Bind(&uf)
	if err != nil {
		c.Error(err).SetMeta("gallery.UpdateController.Bind")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	var gallery m.GalleryType

	// get gallery details
	err = u.Storm.One("ID", uf.ID, &gallery)
	if err != nil {
		c.Error(err).SetMeta("gallery.UpdateController.One")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	gallery.Title = uf.Title
	gallery.Desc = uf.Desc

	// save with updated info
	err = u.Storm.Save(&gallery)
	if err != nil {
		c.Error(err).SetMeta("gallery.UpdateController.Save")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	c.Redirect(http.StatusFound, "/comics/1")

	return

}
