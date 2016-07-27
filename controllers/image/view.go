package controllers

import (
	"encoding/json"
	"net/http"
	"sort"
	"strconv"

	"github.com/boltdb/bolt"
	"github.com/gin-gonic/gin"

	m "fluorescences/models"
	u "fluorescences/utils"
)

// ViewController handles the comic image pages
func ViewController(c *gin.Context) {
	var err error
	var gallery m.GalleryType

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

	// holds our pagination data
	paginate := u.Paged{}
	// holds out page metadata from settings
	metadata, err := u.GetMetadata()
	if err != nil {
		c.Error(err).SetMeta("image.ViewController")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	var image u.FileType
	var title string

	err = u.Bolt.View(func(tx *bolt.Tx) (err error) {
		// the blog bucket
		b := tx.Bucket([]byte(u.GalleryDB))

		cb := b.Cursor()

		_, v := cb.Seek(u.Itob(comicID))

		err = json.Unmarshal(v, &gallery)
		if err != nil {
			return
		}

		paginate.Path = "/image/" + c.Param("id")
		paginate.CurrentPage = currentPage
		paginate.Total = len(gallery.Files)
		paginate.PerPage = 1
		paginate.Desc()

		sort.Sort(gallery.Files)

		for _, c := range gallery.Files {
			if c.ID == currentPage {
				image = c
			}
		}

		title = gallery.Title

		return
	})
	if err != nil {
		c.Error(err).SetMeta("image.ViewController")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	// values for template
	vals := struct {
		Meta  u.Metadata
		Paged u.Paged
		Comic int
		Title string
		Image u.FileType
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
