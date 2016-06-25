package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorController handles error responses
func ErrorController(c *gin.Context) {

	c.HTML(http.StatusOK, "error.tmpl", nil)

	return

}
