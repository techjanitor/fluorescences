package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	m "fluorescences/models"
	u "fluorescences/utils"
)

// NewController posts new galleries
func NewController(c *gin.Context) {

	// holds out page metadata from settings
	metadata, err := u.GetMetadata()
	if err != nil {
		c.Error(err).SetMeta("gallery.NewController")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	vals := struct {
		Meta m.Metadata
	}{
		Meta: metadata,
	}

	c.HTML(http.StatusOK, "gallerynew.tmpl", vals)

	return

}
