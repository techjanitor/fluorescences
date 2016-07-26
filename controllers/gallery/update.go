package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/boltdb/bolt"
	"github.com/gin-gonic/gin"

	m "fluorescences/models"
	u "fluorescences/utils"
)

// UpdateForm is the input from the blog form
type UpdateForm struct {
	ID    int    `form:"id" binding:"required"`
	Title string `form:"title" binding:"required"`
	Desc  string `form:"desc" binding:"required"`
}

// UpdateController updates gallery information
func UpdateController(c *gin.Context) {
	var err error
	var uf UpdateForm

	err = c.Bind(&uf)
	if err != nil {
		c.Error(err).SetMeta("GalleryUpdateController.Bind")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	blog := &m.GalleryType{
		ID:         uf.ID,
		StoredTime: time.Now(),
		Title:      uf.Title,
		Desc:       uf.Desc,
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
func UpdateGallery(blog *m.GalleryType) (err error) {

	// put the tumble in the database
	err = u.Bolt.Update(func(tx *bolt.Tx) (err error) {
		b := tx.Bucket([]byte(u.GalleryDB))

		id := u.Itob(blog.ID)

		cb := b.Cursor()

		_, v := cb.Seek(id)

		var gallery m.GalleryType

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
