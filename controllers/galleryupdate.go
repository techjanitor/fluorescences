package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/boltdb/bolt"
	"github.com/gin-gonic/gin"

	u "fluorescences/utils"
)

// GalleryUpdateForm is the input from the blog form
type GalleryUpdateForm struct {
	ID    int    `form:"id" binding:"required"`
	Title string `form:"title" binding:"required"`
	Desc  string `form:"desc" binding:"required"`
}

// GalleryUpdateController posts new blogs
func GalleryUpdateController(c *gin.Context) {
	var err error
	var guf GalleryUpdateForm

	err = c.Bind(&guf)
	if err != nil {
		c.Error(err).SetMeta("GalleryUpdateController.Bind")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	blog := &GalleryType{
		ID:         guf.ID,
		StoredTime: time.Now(),
		Title:      guf.Title,
		Desc:       guf.Desc,
	}

	err = UpdateGallery(blog)
	if err != nil {
		c.Error(err).SetMeta("GalleryUpdateController.UpdateGallery")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	c.Redirect(http.StatusFound, "/comics/1")

	return

}

// UpdateGallery will add a blog post
func UpdateGallery(blog *GalleryType) (err error) {

	// put the tumble in the database
	err = u.Bolt.Update(func(tx *bolt.Tx) (err error) {
		b := tx.Bucket([]byte(u.GalleryDB))

		id := u.Itob(blog.ID)

		cb := b.Cursor()

		_, v := cb.Seek(id)

		var gallery GalleryType

		err = json.Unmarshal(v, &gallery)
		if err != nil {
			return
		}

		gallery.Title = blog.Title
		gallery.Desc = blog.Desc

		// encode our roomconfig
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
