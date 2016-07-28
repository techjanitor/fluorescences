package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	m "fluorescences/models"
	u "fluorescences/utils"
)

// LoginController is the login page
func LoginController(c *gin.Context) {

	// holds our page metadata from settings
	metadata, err := u.GetMetadata()
	if err != nil {
		c.Error(err).SetMeta("admin.LoginController.GetMetadata")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	vals := struct {
		Meta m.Metadata
	}{
		Meta: metadata,
	}

	c.HTML(http.StatusOK, "login.tmpl", vals)

	return

}
