package controllers

import (
	"net/http"
	"sort"
	"strconv"

	"github.com/gin-gonic/gin"

	m "fluorescences/models"
	u "fluorescences/utils"
)

// ViewController handles the gallery pages
func ViewController(c *gin.Context) {
	var err error
	var gallery m.GalleryType

	comicID, _ := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Error(err).SetMeta("gallery.ViewController")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	currentPage, _ := strconv.Atoi(c.Param("page"))
	if currentPage < 1 {
		currentPage = 1
	}

	// holds our pagination data
	paginate := u.Paged{}
	// holds out page metadata from settings
	metadata, err := u.GetMetadata()
	if err != nil {
		c.Error(err).SetMeta("gallery.ViewController.GetMetadata")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	// get gallery info
	err = u.Storm.One("ID", comicID, &gallery)
	if err != nil {
		c.Error(err).SetMeta("gallery.ViewController")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	sort.Sort(gallery.Files)

	paginate.Path = "/comic/" + c.Param("id")
	paginate.CurrentPage = currentPage
	paginate.Total = len(gallery.Files)
	paginate.PerPage = 9
	paginate.Asc()

	// page through the files slice
	gallery.Files = pageFiles(gallery.Files, paginate.PerPage, paginate.Skip)

	// values for template
	vals := struct {
		Meta    m.Metadata
		Paged   u.Paged
		Gallery m.GalleryType
		All     bool
	}{
		Meta:    metadata,
		Paged:   paginate,
		Gallery: gallery,
		All:     false,
	}

	c.HTML(http.StatusOK, "gallery.tmpl", vals)

	return

}

// pageFiles will page through a files slice with a limit
func pageFiles(files m.Files, limit, skip int) m.Files {
	if skip > len(files) {
		skip = len(files)
	}

	end := skip + limit
	if end > len(files) {
		end = len(files)
	}

	return files[skip:end]
}
