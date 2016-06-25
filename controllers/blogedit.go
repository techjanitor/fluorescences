package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// BlogEditController posts new blogs
func BlogEditController(c *gin.Context) {

	c.HTML(http.StatusOK, "blogedit.tmpl", nil)

	return

}
