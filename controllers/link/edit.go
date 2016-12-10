package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	m "fluorescences/models"
	u "fluorescences/utils"
)

// EditController edits link details
func EditController(c *gin.Context) {
	var err error

	linkID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Error(err).SetMeta("link.EditController")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	// holds out page metadata from settings
	metadata, err := u.GetMetadata()
	if err != nil {
		c.Error(err).SetMeta("link.EditController.GetMetadata")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	var link m.LinkType

	// get the gallery from bolt
	err = u.Storm.One("ID", linkID, &link)
	if err != nil {
		c.Error(err).SetMeta("link.EditController.One")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	vals := struct {
		Meta m.Metadata
		Csrf string
		Link m.LinkType
	}{
		Meta: metadata,
		Csrf: c.MustGet("csrf_token").(string),
		Link: link,
	}

	c.HTML(http.StatusOK, "linkedit.tmpl", vals)

	return

}
