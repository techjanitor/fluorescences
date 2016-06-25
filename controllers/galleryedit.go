package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GalleryEditController posts new blogs
func GalleryEditController(c *gin.Context) {

	c.HTML(http.StatusOK, "galleryedit.tmpl", nil)

	return

}
