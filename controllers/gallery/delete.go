package controllers

import (
	"net/http"

	"github.com/boltdb/bolt"
	"github.com/gin-gonic/gin"

	u "fluorescences/utils"
)

// DeleteForm takes the input from the delete form
type DeleteForm struct {
	Gallery int `form:"gallery" binding:"required"`
}

// DeleteController deletes images from galleries
func DeleteController(c *gin.Context) {
	var err error
	var df DeleteForm

	err = c.Bind(&df)
	if err != nil {
		c.Error(err).SetMeta("gallery.DeleteController.Bind")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	err = DeleteGallery(df.Gallery)
	if err != nil {
		c.Error(err).SetMeta("gallery.DeleteController.DeleteGallery")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	c.Redirect(http.StatusFound, c.Request.Referer())

	return

}

// DeleteGallery will delete a gallery
func DeleteGallery(gallery int) (err error) {

	err = u.Bolt.Update(func(tx *bolt.Tx) (err error) {
		b := tx.Bucket([]byte(u.GalleryDB))

		id := u.Itob(gallery)

		b.Delete(id)

		return

	})
	if err != nil {
		return
	}

	return

}
