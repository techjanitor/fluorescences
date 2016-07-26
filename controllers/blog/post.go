package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/boltdb/bolt"
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
		c.Error(err).SetMeta("BlogPostController")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	blog := m.BlogType{
		User:       "test",
		StoredTime: time.Now(),
		Title:      nf.Title,
		Content:    nf.Post,
	}

	err = AddBlog(blog)
	if err != nil {
		c.Error(err).SetMeta("BlogPostController.AddBlog")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	c.Redirect(http.StatusFound, "/")

	return

}

// AddBlog will add a blog post
func AddBlog(blog m.BlogType) (err error) {

	// put the tumble in the database
	err = u.Bolt.Update(func(tx *bolt.Tx) (err error) {
		bucket := tx.Bucket([]byte(u.BlogDB))

		// get a sequence number
		id, _ := bucket.NextSequence()

		blog.ID = int(id)

		// encode our roomconfig
		encoded, err := json.Marshal(blog)
		if err != nil {
			return
		}

		// put the blog post
		err = bucket.Put(u.Itob(blog.ID), encoded)
		if err != nil {
			return
		}

		return

	})
	if err != nil {
		return
	}

	return

}
