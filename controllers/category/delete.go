package controllers

import (
	"fmt"
	"net/http"

	"github.com/asdine/storm"
	"github.com/gin-gonic/gin"

	m "fluorescences/models"
	u "fluorescences/utils"
)

type deleteForm struct {
	ID int `form:"category" binding:"required"`
}

// DeleteController deletes galleries
func DeleteController(c *gin.Context) {
	var err error
	var df deleteForm

	err = c.Bind(&df)
	if err != nil {
		c.Error(err).SetMeta("category.DeleteController.Bind")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	var galleries []m.GalleryType

	err = u.Storm.Find("Category", df.ID, &galleries)
	if err != nil && err != storm.ErrNotFound {
		c.Error(err).SetMeta("category.DeleteController.Find")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	if len(galleries) > 0 {
		c.Error(fmt.Errorf("Cant delete a category with galleries!")).SetMeta("category.DeleteController.Find")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	var category m.CategoryType

	// get the category from bolt
	err = u.Storm.One("ID", df.ID, &category)
	if err != nil {
		c.Error(err).SetMeta("category.DeleteController.One")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	// delete it
	err = u.Storm.Remove(&category)
	if err != nil {
		c.Error(err).SetMeta("category.DeleteController.Remove")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	c.Redirect(http.StatusFound, c.Request.Referer())

	return

}
