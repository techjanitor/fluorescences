package controllers

import (
	"net/http"
	"sort"
	"strconv"

	"github.com/gin-gonic/gin"

	m "fluorescences/models"
	u "fluorescences/utils"
)

// ViewController handles the comic image pages
func ViewController(c *gin.Context) {
	var err error

	comicID, _ := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Error(err).SetMeta("image.ViewController")
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
		c.Error(err).SetMeta("image.ViewController.GetMetadata")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	var gallery m.GalleryType
	var image m.FileType
	var title string

	err = u.Storm.One("ID", comicID, &gallery)
	if err != nil {
		c.Error(err).SetMeta("image.ViewController.One")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	// holds our pagination data
	paginate := u.Paged{}

	paginate.Path = "/image/" + c.Param("id")
	paginate.CurrentPage = currentPage
	paginate.Total = len(gallery.Files)
	paginate.PerPage = 1
	paginate.Desc()

	sort.Sort(gallery.Files)

	for i, c := range gallery.Files {
		i++
		if i == currentPage {
			image = c
		}
	}

	title = gallery.Title

	// values for template
	vals := struct {
		Meta  m.Metadata
		Paged u.Paged
		Comic int
		Title string
		Image m.FileType
	}{
		Meta:  metadata,
		Paged: paginate,
		Comic: comicID,
		Title: title,
		Image: image,
	}

	c.HTML(http.StatusOK, "image.tmpl", vals)

	return

}
