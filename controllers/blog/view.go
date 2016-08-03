package controllers

import (
	"net/http"
	"strconv"

	"github.com/asdine/storm"
	"github.com/gin-gonic/gin"

	m "fluorescences/models"
	u "fluorescences/utils"
)

// ViewController handles the blog index page
func ViewController(c *gin.Context) {
	var err error

	currentPage, _ := strconv.Atoi(c.Param("page"))
	if currentPage < 1 {
		currentPage = 1
	}

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

	// holds our pagination data
	paginate := u.Paged{}

	paginate.Path = "/blog"
	paginate.CurrentPage = currentPage
	paginate.Total = total
	paginate.PerPage = 10
	paginate.Desc()

	var posts m.Blogs

	// get all the blog posts with a limit
	err = u.Storm.All(&posts, storm.Limit(paginate.PerPage), storm.Skip(paginate.Skip), storm.Reverse())
	if err != nil {
		c.Error(err).SetMeta("blog.ViewController.All")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	for _, post := range posts {
		// format post with markdown
		post.ContentOut = u.Markdown(post.Content)
		// convert time
		post.HumanTime = post.StoredTime.Format("2006-01-02")
	}

	// values for template
	vals := struct {
		Meta  m.Metadata
		Paged u.Paged
		Posts m.Blogs
	}{
		Meta:  metadata,
		Paged: paginate,
		Posts: posts,
	}

	c.HTML(http.StatusOK, "blog.tmpl", vals)

	return

}
