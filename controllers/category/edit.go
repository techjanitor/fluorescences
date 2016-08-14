package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	m "fluorescences/models"
	u "fluorescences/utils"
)

// EditController edits category details
func EditController(c *gin.Context) {
	var err error

	comicID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Error(err).SetMeta("category.EditController")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	// holds out page metadata from settings
	metadata, err := u.GetMetadata()
	if err != nil {
		c.Error(err).SetMeta("category.EditController.GetMetadata")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	var cat m.CategoryType

	// get the gallery from bolt
	err = u.Storm.One("ID", comicID, &cat)
	if err != nil {
		c.Error(err).SetMeta("category.EditController.One")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	vals := struct {
		Meta     m.Metadata
		Csrf     string
		Category m.CategoryType
	}{
		Meta:     metadata,
		Csrf:     c.MustGet("csrf_token").(string),
		Category: cat,
	}

	c.HTML(http.StatusOK, "categoryedit.tmpl", vals)

	return

}
