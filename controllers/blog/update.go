package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	m "fluorescences/models"
	u "fluorescences/utils"
)

type updateForm struct {
	ID    int    `form:"id" binding:"required"`
	Title string `form:"title" binding:"required"`
	Post  string `form:"post" binding:"required"`
}

// UpdateController updates a blog post
func UpdateController(c *gin.Context) {
	var err error
	var uf updateForm

	err = c.Bind(&uf)
	if err != nil {
		c.Error(err).SetMeta("blog.UpdateController.Bind")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	// start transaction
	tx, err := u.Storm.Begin(true)
	if err != nil {
		c.Error(err).SetMeta("blog.UpdateController.Begin")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}
	defer tx.Rollback()

	var blog m.BlogType

	// get category details
	err = tx.One("ID", uf.ID, &blog)
	if err != nil {
		c.Error(err).SetMeta("blog.UpdateController.One")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	blog.Title = uf.Title
	blog.Content = uf.Post

	// save with updated info
	err = tx.Save(&blog)
	if err != nil {
		c.Error(err).SetMeta("blog.UpdateController.Save")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	// commit
	tx.Commit()

	c.Redirect(http.StatusFound, "/")

	return

}
