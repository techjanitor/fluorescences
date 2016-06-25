package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	u "fluorescences/utils"
)

// GalleryNewController posts new blogs
func GalleryNewController(c *gin.Context) {

	// holds out page metadata from settings
	metadata, err := u.GetMetadata()
	if err != nil {
		c.Error(err).SetMeta("GalleryNewController")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	vals := struct {
		Meta u.Metadata
		New  bool
		Edit bool
	}{
		Meta: metadata,
		New:  true,
		Edit: false,
	}

	c.HTML(http.StatusOK, "galleryedit.tmpl", vals)

	return

}
