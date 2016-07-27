package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	m "fluorescences/models"
	u "fluorescences/utils"
)

// EditController edits gallery details
func EditController(c *gin.Context) {
	var err error
	var gallery m.GalleryType

	comicID, _ := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Error(err).SetMeta("gallery.EditController")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	// holds out page metadata from settings
	metadata, err := u.GetMetadata()
	if err != nil {
		c.Error(err).SetMeta("gallery.EditController")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	err = u.Storm.One("ID", comicID, &gallery)
	if err != nil {
		c.Error(err).SetMeta("gallery.EditController.Storm")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	vals := struct {
		Meta    m.Metadata
		Gallery m.GalleryType
	}{
		Meta:    metadata,
		Gallery: gallery,
	}

	c.HTML(http.StatusOK, "galleryedit.tmpl", vals)

	return

}
