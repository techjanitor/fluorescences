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

	comicID, _ := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Error(err).SetMeta("gallery.EditController")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	// holds out page metadata from settings
	metadata, err := u.GetMetadata()
	if err != nil {
		c.Error(err).SetMeta("gallery.EditController.GetMetadata")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	var gallery m.GalleryType

	// get the gallery from bolt
	err = u.Storm.One("ID", comicID, &gallery)
	if err != nil {
		c.Error(err).SetMeta("gallery.EditController.One")
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
		Gallery    m.GalleryType
		Categories m.Categories
	}{
		Meta:       metadata,
		Csrf:       c.MustGet("csrf_token").(string),
		Gallery:    gallery,
		Categories: cats,
	}

	c.HTML(http.StatusOK, "galleryedit.tmpl", vals)

	return

}
