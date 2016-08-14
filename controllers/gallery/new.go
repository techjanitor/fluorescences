package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	m "fluorescences/models"
	u "fluorescences/utils"
)

// NewController is the handler for the page where you create a new gallery
func NewController(c *gin.Context) {

	// holds out page metadata from settings
	metadata, err := u.GetMetadata()
	if err != nil {
		c.Error(err).SetMeta("gallery.NewController.GetMetadata")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	var cats m.Categories

	// get all the categories
	err = u.Storm.All(&cats)
	if err != nil {
		c.Error(err).SetMeta("gallery.NewController.All")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	vals := struct {
		Meta       m.Metadata
		Csrf       string
		Categories m.Categories
	}{
		Meta:       metadata,
		Csrf:       c.MustGet("csrf_token").(string),
		Categories: cats,
	}

	c.HTML(http.StatusOK, "gallerynew.tmpl", vals)

	return

}
