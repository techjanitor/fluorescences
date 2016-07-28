package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	u "fluorescences/utils"
)

type loginForm struct {
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginController handles the admin gallery page
func LoginController(c *gin.Context) {
	var err error
	var lf loginForm

	err = c.Bind(&lf)
	if err != nil {
		c.Error(err).SetMeta("admin.LoginController.Bind")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	// compare passwords
	err = u.CheckPassword(lf.Password)
	if err != nil {
		c.Error(err).SetMeta("admin.LoginController.ComparePassword")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	// create jwt token
	token, err := u.MakeToken()
	if err != nil {
		c.Error(err).SetMeta("admin.LoginController.MakeToken")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	// set the jwt cookie
	http.SetCookie(c.Writer, u.CreateCookie(token))

	c.Redirect(http.StatusFound, "/admin/panel")

	return

}
