package controllers

import (
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"

	m "fluorescences/models"
	u "fluorescences/utils"
)

type updateForm struct {
	ID    int    `form:"id" binding:"required"`
	Title string `form:"title" binding:"required"`
	URL   string `form:"url" binding:"required"`
}

// UpdateController updates image information
func UpdateController(c *gin.Context) {
	var err error
	var uf updateForm

	err = c.Bind(&uf)
	if err != nil {
		c.Error(err).SetMeta("image.UpdateController.Bind")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	// start transaction
	tx, err := u.Storm.Begin(true)
	if err != nil {
		c.Error(err).SetMeta("blog.UpdateController.Begin")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}
	defer tx.Rollback()

	var link m.LinkType

	// get category details
	err = tx.One("ID", uf.ID, &link)
	if err != nil {
		c.Error(err).SetMeta("link.UpdateController.One")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	address, err := url.Parse(uf.URL)
	if err != nil {
		c.Error(err).SetMeta("link.UpdateController.Parse")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	link.Title = uf.Title
	link.Address = address.String()
	link.URL = address

	// save with updated info
	err = tx.Save(&link)
	if err != nil {
		c.Error(err).SetMeta("link.UpdateController.Save")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	// commit
	tx.Commit()

	c.Redirect(http.StatusFound, "/admin/panel")

	return

}
