package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	m "fluorescences/models"
	u "fluorescences/utils"
)

// NewForm is the input from the blog form
type NewForm struct {
	Title string `form:"title" binding:"required"`
	Post  string `form:"post" binding:"required"`
}

// PostController posts new blogs
func PostController(c *gin.Context) {
	var err error
	var nf NewForm

	err = c.Bind(&nf)
	if err != nil {
		c.Error(err).SetMeta("blog.PostController.Bind")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	blog := m.BlogType{
		User:       "test",
		StoredTime: time.Now(),
		Title:      nf.Title,
		Content:    nf.Post,
	}

	err = u.Storm.Save(&blog)
	if err != nil {
		c.Error(err).SetMeta("blog.PostController.Save")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	c.Redirect(http.StatusFound, "/")

	return

}
