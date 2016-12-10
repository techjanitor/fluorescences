package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	u "fluorescences/utils"
)

type passwordForm struct {
	Old   string `form:"oldpassword" binding:"required"`
	New   string `form:"newpassword" binding:"required"`
	Check string `form:"checkpassword" binding:"required"`
}

// PasswordController updates the password
func PasswordController(c *gin.Context) {
	var err error
	var pf passwordForm

	err = c.Bind(&pf)
	if err != nil {
		c.Error(err).SetMeta("admin.PasswordController.Bind")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	err = u.CheckPassword(pf.Old)
	if err != nil {
		c.Error(err).SetMeta("admin.PasswordController.Check")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	if pf.New != pf.Check {
		c.Error(err).SetMeta("admin.PasswordController.Compare")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	var hash []byte

	hash, err = u.HashPassword(pf.Check)
	if err != nil {
		c.Error(err).SetMeta("admin.PasswordController.HashPassword")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	var user u.User

	err = u.Storm.Get("auth", "user", &user)
	if err != nil {
		c.Error(err).SetMeta("admin.PasswordController.HashPassword")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	// set user password
	user.Password = hash

	err = u.Storm.Set("auth", "user", &user)
	if err != nil {
		c.Error(err).SetMeta("admin.PasswordController.Set")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	// unset the jwt cookie
	http.SetCookie(c.Writer, u.DeleteCookie())

	c.Redirect(http.StatusFound, "/admin/panel")

	return

}
