package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	u "fluorescences/utils"
)

type loginForm struct {
	Name     string `form:"name" binding:"required"`
	Password string `form:"password" binding:"required"`
}

// AuthController handles user credential authentication
func AuthController(c *gin.Context) {
	var err error
	var lf loginForm

	err = c.Bind(&lf)
	if err != nil {
		c.Error(err).SetMeta("admin.AuthController.Bind")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	// compare passwords
	err = u.CheckPassword(lf.Password)
	if err != nil {
		c.Error(err).SetMeta("admin.AuthController.ComparePassword")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	// create jwt token
	token, err := u.MakeToken()
	if err != nil {
		c.Error(err).SetMeta("admin.AuthController.MakeToken")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	// set the jwt cookie
	http.SetCookie(c.Writer, u.CreateCookie(token))

	c.Redirect(http.StatusFound, "/admin/panel")

	return

}
