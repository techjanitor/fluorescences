package controllers

import (
	"net/http"

	"github.com/asdine/storm"
	"github.com/gin-gonic/gin"

	m "fluorescences/models"
	u "fluorescences/utils"
)

// GalleryController handles the admin gallery page
func GalleryController(c *gin.Context) {
	var err error
	var galleries m.Galleries

	// holds out page metadata from settings
	metadata, err := u.GetMetadata()
	if err != nil {
		c.Error(err).SetMeta("admin.GalleryController")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	u.Storm.All(&galleries, storm.Reverse())

	// values for template
	vals := struct {
		Meta      m.Metadata
		Csrf      string
		Galleries m.Galleries
	}{
		Meta:      metadata,
		Csrf:      c.MustGet("csrf_token").(string),
		Galleries: galleries,
	}

	c.HTML(http.StatusOK, "panel.tmpl", vals)

	return

}
