package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	m "fluorescences/models"
	u "fluorescences/utils"
)

type privateForm struct {
	ID      int  `form:"id" binding:"required"`
	Private bool `form:"private"`
}

// PrivateController updates a galleries private status
func PrivateController(c *gin.Context) {
	var err error
	var pf privateForm

	err = c.Bind(&pf)
	if err != nil {
		c.Error(err).SetMeta("gallery.PrivateController.Bind")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	// start transaction
	tx, err := u.Storm.Begin(true)
	if err != nil {
		c.Error(err).SetMeta("gallery.PrivateController.Begin")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}
	defer tx.Rollback()

	var gallery m.GalleryType

	// get gallery details
	err = tx.One("ID", pf.ID, &gallery)
	if err != nil {
		c.Error(err).SetMeta("gallery.PrivateController.One")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	gallery.Private = pf.Private

	// save with updated info
	err = tx.Save(&gallery)
	if err != nil {
		c.Error(err).SetMeta("gallery.PrivateController.Save")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	// commit
	tx.Commit()

	c.Redirect(http.StatusFound, "/admin/panel")

	return

}
