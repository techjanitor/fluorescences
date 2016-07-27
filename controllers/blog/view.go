package controllers

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/asdine/storm"
	"github.com/gin-gonic/gin"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"

	m "fluorescences/models"
	u "fluorescences/utils"
)

// ViewController handles the blog index page
func ViewController(c *gin.Context) {
	var err error
	var posts []*m.BlogType

	currentPage, _ := strconv.Atoi(c.Param("page"))
	if currentPage < 1 {
		currentPage = 1
	}

	// holds our pagination data
	paginate := u.Paged{}

	// holds out page metadata from settings
	metadata, err := u.GetMetadata()
	if err != nil {
		c.Error(err).SetMeta("blog.ViewController.GetMetadata")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	// get a count of the blogs
	total, err := u.Storm.Count(&m.BlogType{})
	if err != nil {
		c.Error(err).SetMeta("blog.ViewController.Count")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	paginate.Path = "/blog"
	paginate.CurrentPage = currentPage
	paginate.Total = total
	paginate.PerPage = 10
	paginate.Desc()

	// get all the blog posts with a limit
	err = u.Storm.All(&posts, storm.Limit(paginate.PerPage), storm.Skip(paginate.Skip), storm.Reverse())
	if err != nil {
		c.Error(err).SetMeta("blog.ViewController.All")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	for _, post := range posts {
		// make the post formatted with markdown
		unsafe := blackfriday.MarkdownCommon([]byte(post.Content))
		// sanitize the input
		html := bluemonday.UGCPolicy().SanitizeBytes(unsafe)
		// convert to template format
		post.ContentOut = template.HTML(html)
		// convert time
		post.HumanTime = post.StoredTime.Format("2006-01-02")
	}

	// values for template
	vals := struct {
		Meta  m.Metadata
		Paged u.Paged
		Posts []*m.BlogType
	}{
		Meta:  metadata,
		Paged: paginate,
		Posts: posts,
	}

	c.HTML(http.StatusOK, "blog.tmpl", vals)

	return

}
