package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	m "fluorescences/models"
	u "fluorescences/utils"
)

// NewController is the handler for the page where you create a new category
func NewController(c *gin.Context) {

	// holds out page metadata from settings
	metadata, err := u.GetMetadata()
	if err != nil {
		c.Error(err).SetMeta("category.NewController.GetMetadata")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	vals := struct {
		Meta m.Metadata
		Csrf string
	}{
		Meta: metadata,
		Csrf: c.MustGet("csrf_token").(string),
	}

	c.HTML(http.StatusOK, "categorynew.tmpl", vals)

	return

}
