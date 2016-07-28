package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	m "fluorescences/models"
	u "fluorescences/utils"
)

type deleteForm struct {
	ID int `form:"id" binding:"required"`
}

// DeleteController deletes blog posts
func DeleteController(c *gin.Context) {
	var err error
	var df deleteForm

	err = c.Bind(&df)
	if err != nil {
		c.Error(err).SetMeta("blog.DeleteController.Bind")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	var blog m.BlogType

	// get the gallery from bolt
	err = u.Storm.One("ID", df.ID, &blog)
	if err != nil {
		c.Error(err).SetMeta("blog.DeleteController.One")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	// delete it
	err = u.Storm.Remove(&blog)
	if err != nil {
		c.Error(err).SetMeta("blog.DeleteController.Remove")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	c.Redirect(http.StatusFound, c.Request.Referer())

	return

}
