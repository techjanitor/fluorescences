package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/asdine/storm"
	"github.com/gin-gonic/gin"

	m "fluorescences/models"
	u "fluorescences/utils"
)

// IndexController handles the galleries index page
func IndexController(c *gin.Context) {
	var err error
	var galleries []*m.GalleryType

	currentPage, _ := strconv.Atoi(c.Param("page"))
	if currentPage < 1 {
		currentPage = 1
	}

	// holds our pagination data
	paginate := u.Paged{}

	// holds out page metadata from settings
	metadata, err := u.GetMetadata()
	if err != nil {
		c.Error(err).SetMeta("gallery.IndexController")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	total, err := u.Storm.Count(&m.GalleryType{})
	if err != nil {
		c.Error(err).SetMeta("blog.ViewController.Storm")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	paginate.Path = "/comics"
	paginate.CurrentPage = currentPage
	paginate.Total = total
	paginate.PerPage = 5
	paginate.Desc()

	fmt.Println(paginate)

	err = u.Storm.All(&galleries, storm.Limit(paginate.PerPage), storm.Skip(paginate.Skip))
	if err != nil {
		c.Error(err).SetMeta("blog.ViewController.Storm")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	for _, gallery := range galleries {
		// convert time
		gallery.HumanTime = gallery.StoredTime.Format("2006-01-02")
		// cover image is the first image in the slice
		gallery.Cover = gallery.Files[0].Filename
	}

	// values for template
	vals := struct {
		Meta      m.Metadata
		Paged     u.Paged
		Galleries []*m.GalleryType
		All       bool
	}{
		Meta:      metadata,
		Paged:     paginate,
		Galleries: galleries,
		All:       true,
	}

	c.HTML(http.StatusOK, "gallery.tmpl", vals)

	return

}
