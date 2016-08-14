package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	m "fluorescences/models"
	u "fluorescences/utils"
)

type updateForm struct {
	Title string `form:"title" binding:"required"`
	Desc  string `form:"desc" binding:"required"`
}

// UpdateController updates the site settings
func UpdateController(c *gin.Context) {
	var err error
	var uf updateForm

	err = c.Bind(&uf)
	if err != nil {
		c.Error(err).SetMeta("admin.UpdateController.Bind")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	var meta m.Metadata

	meta.Title = uf.Title
	meta.Desc = uf.Desc

	err = u.Storm.Set("data", "metadata", &meta)
	if err != nil {
		return
	}

	c.Redirect(http.StatusFound, "/admin/panel")

	return

}
