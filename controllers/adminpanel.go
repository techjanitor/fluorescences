package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/boltdb/bolt"
	"github.com/gin-gonic/gin"

	u "fluorescences/utils"
)

// AdminPanelController is the main admin menu
func AdminPanelController(c *gin.Context) {
	var err error
	var galleries []GalleryType

	// holds out page metadata from settings
	metadata, err := u.GetMetadata()
	if err != nil {
		c.Error(err).SetMeta("GalleryController")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	err = u.Bolt.View(func(tx *bolt.Tx) (err error) {
		// the blog bucket
		b := tx.Bucket([]byte(GalleryDB))

		cb := b.Cursor()

		for k, v := cb.Last(); k != nil; k, v = cb.Prev() {

			gallery := GalleryType{}

			err = json.Unmarshal(v, &gallery)
			if err != nil {
				return
			}

			// convert time
			gallery.HumanTime = gallery.StoredTime.Format("2006-01-02")

			galleries = append(galleries, gallery)

		}
		return
	})
	if err != nil {
		c.Error(err).SetMeta("GalleryController")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	// values for template
	vals := struct {
		Meta      u.Metadata
		Galleries []GalleryType
	}{
		Meta:      metadata,
		Galleries: galleries,
	}

	c.HTML(http.StatusOK, "panel.tmpl", vals)

	return

}
