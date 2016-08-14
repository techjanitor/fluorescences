package controllers

import (
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

	categoryID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Error(err).SetMeta("gallery.IndexController")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	currentPage, _ := strconv.Atoi(c.Param("page"))
	if currentPage < 1 {
		currentPage = 1
	}

	// holds out page metadata from settings
	metadata, err := u.GetMetadata()
	if err != nil {
		c.Error(err).SetMeta("gallery.IndexController.GetMetadata")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	var total []m.GalleryType

	err = u.Storm.Find("Category", categoryID, &total)
	if err != nil {
		c.Error(err).SetMeta("gallery.IndexController.Find")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	// holds our pagination data
	paginate := u.Paged{}

	paginate.Path = "/comics"
	paginate.CurrentPage = currentPage
	paginate.Total = len(total)
	paginate.PerPage = 6
	paginate.Desc()

	var galleries []m.GalleryType

	// get all the galleries with a limit
	err = u.Storm.Find("Category", categoryID, &galleries, storm.Limit(paginate.PerPage), storm.Skip(paginate.Skip), storm.Reverse())
	if err != nil {
		c.Error(err).SetMeta("gallery.IndexController.Find")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	for idx := range galleries {
		// cover image is the first image in the slice
		galleries[idx].Cover = galleries[idx].Files[0].Filename
	}

	// values for template
	vals := struct {
		Meta      m.Metadata
		Paged     u.Paged
		Galleries []m.GalleryType
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
