package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	m "fluorescences/models"
	u "fluorescences/utils"
)

// InputController is the handler for the page where people can enter the gallery password
func InputController(c *gin.Context) {

	comicID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Error(err).SetMeta("keys.InputController")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	// holds out page metadata from settings
	metadata, err := u.GetMetadata()
	if err != nil {
		c.Error(err).SetMeta("keys.InputController.GetMetadata")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	vals := struct {
		Meta m.Metadata
		ID   int
	}{
		Meta: metadata,
		ID:   comicID,
	}

	c.HTML(http.StatusOK, "gallerypassword.tmpl", vals)

	return

}
