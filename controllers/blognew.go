package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	u "fluorescences/utils"
)

// BlogNewController posts new blogs
func BlogNewController(c *gin.Context) {

	// holds out page metadata from settings
	metadata, err := u.GetMetadata()
	if err != nil {
		c.Error(err).SetMeta("BlogNewController")
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

	c.HTML(http.StatusOK, "blogedit.tmpl", vals)

	return

}
