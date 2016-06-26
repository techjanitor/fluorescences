package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/boltdb/bolt"
	"github.com/gin-gonic/gin"

	u "fluorescences/utils"
)

// GalleryEditController edits gallery details
func GalleryEditController(c *gin.Context) {
	var err error
	var gallery GalleryType

	comicID, _ := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Error(err).SetMeta("ComicController")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	// holds out page metadata from settings
	metadata, err := u.GetMetadata()
	if err != nil {
		c.Error(err).SetMeta("GalleryEditController")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	err = u.Bolt.View(func(tx *bolt.Tx) (err error) {
		// the blog bucket
		b := tx.Bucket([]byte(u.GalleryDB))

		cb := b.Cursor()

		_, v := cb.Seek(u.Itob(comicID))

		err = json.Unmarshal(v, &gallery)
		if err != nil {
			return
		}

		return
	})
	if err != nil {
		c.Error(err).SetMeta("ComicController")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	vals := struct {
		Meta    u.Metadata
		Gallery GalleryType
	}{
		Meta:    metadata,
		Gallery: gallery,
	}

	c.HTML(http.StatusOK, "galleryedit.tmpl", vals)

	return

}
