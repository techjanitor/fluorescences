package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	u "fluorescences/utils"
)

// LogoutController will delete the login cookie
func LogoutController(c *gin.Context) {

	// unset the jwt cookie
	http.SetCookie(c.Writer, u.DeleteCookie())

	c.Redirect(http.StatusFound, "/")

	return

}
