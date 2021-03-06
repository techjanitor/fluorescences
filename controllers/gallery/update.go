package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	m "fluorescences/models"
	u "fluorescences/utils"
)

type updateForm struct {
	ID       int    `form:"id" binding:"required"`
	Category int    `form:"category" binding:"required"`
	Title    string `form:"title" binding:"required"`
	Desc     string `form:"desc" binding:"required"`
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

	// start transaction
	tx, err := u.Storm.Begin(true)
	if err != nil {
		c.Error(err).SetMeta("gallery.UpdateController.Begin")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}
	defer tx.Rollback()

	var gallery m.GalleryType

	// get gallery details
	err = tx.One("ID", uf.ID, &gallery)
	if err != nil {
		c.Error(err).SetMeta("gallery.UpdateController.One")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	gallery.Title = uf.Title
	gallery.Desc = uf.Desc
	gallery.Category = uf.Category

	// save with updated info
	err = tx.Save(&gallery)
	if err != nil {
		c.Error(err).SetMeta("gallery.UpdateController.Save")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	// commit
	tx.Commit()

	c.Redirect(http.StatusFound, fmt.Sprintf("/comics/%d/1", uf.Category))

	return

}
