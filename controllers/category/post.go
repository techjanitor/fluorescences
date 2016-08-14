package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	m "fluorescences/models"
	u "fluorescences/utils"
)

type newForm struct {
	Title string `form:"title" binding:"required"`
	Desc  string `form:"desc" binding:"required"`
}

// PostController posts new categories
func PostController(c *gin.Context) {
	var err error
	var nf newForm

	err = c.Bind(&nf)
	if err != nil {
		c.Error(err).SetMeta("category.PostController.Bind")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	category := m.CategoryType{
		Title: nf.Title,
		Desc:  nf.Desc,
	}

	// save category
	err = u.Storm.Save(&category)
	if err != nil {
		c.Error(err).SetMeta("category.PostController.Save")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	c.Redirect(http.StatusFound, "/admin/panel")

	return

}
