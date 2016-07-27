package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/boltdb/bolt"
	"github.com/gin-gonic/gin"

	m "fluorescences/models"
	u "fluorescences/utils"
)

// IndexController handles the galleries index page
func IndexController(c *gin.Context) {
	var err error
	var galleries []m.GalleryType

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

	err = u.Bolt.View(func(tx *bolt.Tx) (err error) {
		// the blog bucket
		b := tx.Bucket([]byte(u.GalleryDB))

		// stats for key count
		stats := b.Stats()

		fmt.Println(stats.KeyN)

		paginate.Path = "/comics"
		paginate.CurrentPage = currentPage
		paginate.Total = stats.KeyN
		paginate.PerPage = 10
		paginate.Desc()

		cb := b.Cursor()

		for k, v := cb.Seek(u.Itob(paginate.Start)); k != nil && !bytes.Equal(k, u.Itob(paginate.End)); k, v = cb.Prev() {

			gallery := m.GalleryType{}

			err = json.Unmarshal(v, &gallery)
			if err != nil {
				return
			}

			// convert time
			gallery.HumanTime = gallery.StoredTime.Format("2006-01-02")

			gallery.Cover = gallery.Files[0].Filename

			galleries = append(galleries, gallery)

		}
		return
	})
	if err != nil {
		c.Error(err).SetMeta("gallery.IndexController")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	// values for template
	vals := struct {
		Meta      u.Metadata
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
