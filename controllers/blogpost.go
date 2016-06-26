package controllers

import (
	"encoding/json"
	"html/template"
	"net/http"
	"time"

	"github.com/boltdb/bolt"
	"github.com/gin-gonic/gin"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"

	u "fluorescences/utils"
)

// BlogType holds a blog post
type BlogType struct {
	ID         int
	User       string
	Title      string
	Content    template.HTML
	HumanTime  string
	StoredTime time.Time
}

// BlogForm is the input from the blog form
type BlogForm struct {
	Title string `form:"title" binding:"required"`
	Post  string `form:"post" binding:"required"`
}

// BlogPostController posts new blogs
func BlogPostController(c *gin.Context) {
	var err error
	var bf BlogForm

	err = c.Bind(&bf)
	if err != nil {
		c.Error(err).SetMeta("BlogPostController")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	// make the post formatted with markdown
	unsafe := blackfriday.MarkdownCommon([]byte(bf.Post))
	// sanitize the input
	html := bluemonday.UGCPolicy().SanitizeBytes(unsafe)
	// convert to template format
	content := template.HTML(html)

	blog := BlogType{
		User:       "test",
		StoredTime: time.Now(),
		Title:      bf.Title,
		Content:    content,
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
func AddBlog(blog BlogType) (err error) {

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
