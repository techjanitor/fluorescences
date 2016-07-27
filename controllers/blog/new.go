package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	m "fluorescences/models"
	u "fluorescences/utils"
)

// NewController posts new blogs
func NewController(c *gin.Context) {

	// holds our page metadata from settings
	metadata, err := u.GetMetadata()
	if err != nil {
		c.Error(err).SetMeta("blog.NewController.GetMetadata")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	vals := struct {
		Meta m.Metadata
		New  bool
		Edit bool
	}{
		Meta: metadata,
		New:  true,
		Edit: false,
	}

	c.HTML(http.StatusOK, "blogedit.tmpl", vals)

	return

}
