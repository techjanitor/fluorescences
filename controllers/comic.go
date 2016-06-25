package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/boltdb/bolt"
	"github.com/gin-gonic/gin"

	u "fluorescences/utils"
)

// ComicController handles the comic gallery pages
func ComicController(c *gin.Context) {
	var err error
	var gallery GalleryType

	comicID, _ := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Error(err).SetMeta("ComicController")
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
		c.Error(err).SetMeta("ComicController")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	err = u.Bolt.View(func(tx *bolt.Tx) (err error) {
		// the blog bucket
		b := tx.Bucket([]byte(GalleryDB))

		cb := b.Cursor()

		_, v := cb.Seek(u.Itob(comicID))

		err = json.Unmarshal(v, &gallery)
		if err != nil {
			return
		}

		paginate.CurrentPage = currentPage
		paginate.Total = len(gallery.Files)
		paginate.PerPage = 20
		paginate.Desc()

		return
	})
	if err != nil {
		c.Error(err).SetMeta("ComicController")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	// values for template
	vals := struct {
		Meta    u.Metadata
		Paged   u.Paged
		Gallery GalleryType
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
