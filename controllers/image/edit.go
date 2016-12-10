package controllers

import (
	"net/http"
	"sort"
	"strconv"

	"github.com/gin-gonic/gin"

	m "fluorescences/models"
	u "fluorescences/utils"
)

// EditController edits category details
func EditController(c *gin.Context) {
	var err error

	galleryID, err := strconv.Atoi(c.Param("gallery"))
	if err != nil {
		c.Error(err).SetMeta("image.EditController")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	imageID, err := strconv.Atoi(c.Param("image"))
	if err != nil {
		c.Error(err).SetMeta("image.EditController")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	// holds out page metadata from settings
	metadata, err := u.GetMetadata()
	if err != nil {
		c.Error(err).SetMeta("image.EditController.GetMetadata")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	var gallery m.GalleryType
	var image m.FileType

	err = u.Storm.One("ID", galleryID, &gallery)
	if err != nil {
		c.Error(err).SetMeta("image.EditController.One")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	sort.Sort(gallery.Files)

	// update the image info
	for i := len(gallery.Files) - 1; i >= 0; i-- {
		file := gallery.Files[i]

		if file.ID == imageID {
			image = file
		}
	}

	vals := struct {
		Meta    m.Metadata
		Csrf    string
		Gallery m.GalleryType
		Image   m.FileType
	}{
		Meta:    metadata,
		Csrf:    c.MustGet("csrf_token").(string),
		Gallery: gallery,
		Image:   image,
	}

	c.HTML(http.StatusOK, "imageedit.tmpl", vals)

	return

}
