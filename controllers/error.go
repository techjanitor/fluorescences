package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorController handles no matching routes
func ErrorController(c *gin.Context) {

	c.HTML(http.StatusNotFound, "error.tmpl", nil)

	return

}
