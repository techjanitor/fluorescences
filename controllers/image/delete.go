package controllers

import (
	"encoding/json"
	"net/http"
	"sort"

	"github.com/boltdb/bolt"
	"github.com/gin-gonic/gin"

	m "fluorescences/models"
	u "fluorescences/utils"
)

// DeleteForm takes the input from the delete form
type DeleteForm struct {
	Gallery int `form:"gallery" binding:"required"`
	Image   int `form:"image" binding:"required"`
}

// DeleteController deletes images from galleries
func DeleteController(c *gin.Context) {
	var err error
	var df DeleteForm

	err = c.Bind(&df)
	if err != nil {
		c.Error(err).SetMeta("image.DeleteController.Bind")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	err = DeleteImage(df.Gallery, df.Image)
	if err != nil {
		c.Error(err).SetMeta("image.DeleteController.DeleteImage")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	c.Redirect(http.StatusFound, c.Request.Referer())

	return

}

// DeleteImage will delete an image
func DeleteImage(gallery, image int) (err error) {

	err = u.Bolt.Update(func(tx *bolt.Tx) (err error) {
		b := tx.Bucket([]byte(u.GalleryDB))

		id := u.Itob(gallery)

		cb := b.Cursor()

		_, v := cb.Seek(id)

		var gallery m.GalleryType

		err = json.Unmarshal(v, &gallery)
		if err != nil {
			return
		}

		sort.Sort(gallery.Files)

		for i := len(gallery.Files) - 1; i >= 0; i-- {
			file := gallery.Files[i]

			if file.ID == image {
				gallery.Files = append(gallery.Files[:i],
					gallery.Files[i+1:]...)
			}
		}

		encoded, err := json.Marshal(gallery)
		if err != nil {
			return
		}

		// put the blog post
		err = b.Put(id, encoded)
		if err != nil {
			return
		}

		return

	})
	if err != nil {
		return
	}

	return

}
