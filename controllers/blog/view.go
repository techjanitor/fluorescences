package controllers

import (
	"bytes"
	"encoding/json"
	"html/template"
	"net/http"
	"strconv"

	"github.com/boltdb/bolt"
	"github.com/gin-gonic/gin"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"

	m "fluorescences/models"
	u "fluorescences/utils"
)

// ViewController handles the blog index page
func ViewController(c *gin.Context) {
	var err error
	var posts []m.BlogType

	currentPage, _ := strconv.Atoi(c.Param("page"))
	if currentPage < 1 {
		currentPage = 1
	}

	// holds our pagination data
	paginate := u.Paged{}
	// holds out page metadata from settings
	metadata, err := u.GetMetadata()
	if err != nil {
		c.Error(err).SetMeta("BlogController.GetMetadata")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	err = u.Bolt.View(func(tx *bolt.Tx) (err error) {
		// the blog bucket
		b := tx.Bucket([]byte(u.BlogDB))

		// stats for key count
		stats := b.Stats()

		paginate.Path = "/blog"
		paginate.CurrentPage = currentPage
		paginate.Total = stats.KeyN
		paginate.PerPage = 10
		paginate.Desc()

		cb := b.Cursor()

		for k, v := cb.Seek(u.Itob(paginate.Start)); k != nil && !bytes.Equal(k, u.Itob(paginate.End)); k, v = cb.Prev() {

			post := m.BlogType{}

			err = json.Unmarshal(v, &post)
			if err != nil {
				return
			}

			// convert time
			post.HumanTime = post.StoredTime.Format("2006-01-02")

			// make the post formatted with markdown
			unsafe := blackfriday.MarkdownCommon([]byte(post.Content))
			// sanitize the input
			html := bluemonday.UGCPolicy().SanitizeBytes(unsafe)
			// convert to template format
			post.ContentOut = template.HTML(html)

			posts = append(posts, post)

		}
		return
	})
	if err != nil {
		c.Error(err).SetMeta("BlogController")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	// values for template
	vals := struct {
		Meta  u.Metadata
		Paged u.Paged
		Posts []m.BlogType
	}{
		Meta:  metadata,
		Paged: paginate,
		Posts: posts,
	}

	c.HTML(http.StatusOK, "blog.tmpl", vals)

	return

}
