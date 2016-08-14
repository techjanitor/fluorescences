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

// UpdateController updates category information
func UpdateController(c *gin.Context) {
	var err error
	var uf updateForm

	err = c.Bind(&uf)
	if err != nil {
		c.Error(err).SetMeta("category.UpdateController.Bind")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	// start transaction
	tx, err := u.Storm.Begin(true)
	if err != nil {
		c.Error(err).SetMeta("category.UpdateController.Begin")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}
	defer tx.Rollback()

	var category m.CategoryType

	// get category details
	err = tx.One("ID", uf.ID, &category)
	if err != nil {
		c.Error(err).SetMeta("category.UpdateController.One")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	category.Title = uf.Title
	category.Desc = uf.Desc

	// save with updated info
	err = tx.Save(&category)
	if err != nil {
		c.Error(err).SetMeta("category.UpdateController.Save")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	// commit
	tx.Commit()

	c.Redirect(http.StatusFound, "/admin/panel")

	return

}
