package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	m "fluorescences/models"
	u "fluorescences/utils"
)

// EditController edits category details
func EditController(c *gin.Context) {
	var err error

	blogID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Error(err).SetMeta("blog.EditController")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	// holds out page metadata from settings
	metadata, err := u.GetMetadata()
	if err != nil {
		c.Error(err).SetMeta("blog.EditController.GetMetadata")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	var blog m.BlogType

	// get the gallery from bolt
	err = u.Storm.One("ID", blogID, &blog)
	if err != nil {
		c.Error(err).SetMeta("blog.EditController.One")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	vals := struct {
		Meta m.Metadata
		Csrf string
		Blog m.BlogType
		New  bool
		Edit bool
	}{
		Meta: metadata,
		Csrf: c.MustGet("csrf_token").(string),
		Blog: blog,
		New:  false,
		Edit: true,
	}

	c.HTML(http.StatusOK, "blogedit.tmpl", vals)

	return

}
