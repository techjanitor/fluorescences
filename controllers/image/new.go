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

// NewForm is the input from the new image form
type NewForm struct {
	ID int `form:"id" binding:"required"`
}

// NewController posts new blogs
func NewController(c *gin.Context) {
	var err error
	var inf NewForm

	err = c.Bind(&inf)
	if err != nil {
		c.Error(err).SetMeta("image.NewController.Bind")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	// Check if theres a file
	upload, fileheader, err := c.Request.FormFile("file")
	if err != nil {
		c.Error(err).SetMeta("image.NewController")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	file, err := u.SaveFile(upload, fileheader.Filename)
	if err != nil {
		c.Error(err).SetMeta("image.NewController.SaveFile")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	err = AddImage(inf.ID, file)
	if err != nil {
		c.Error(err).SetMeta("image.NewController.AddImage")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	c.Redirect(http.StatusFound, c.Request.Referer())

	return

}

// AddImage will add a blog post
func AddImage(gid int, file u.FileType) (err error) {

	err = u.Bolt.Update(func(tx *bolt.Tx) (err error) {
		b := tx.Bucket([]byte(u.GalleryDB))

		id := u.Itob(gid)

		cb := b.Cursor()

		_, v := cb.Seek(id)

		var gallery m.GalleryType

		err = json.Unmarshal(v, &gallery)
		if err != nil {
			return
		}

		sort.Sort(gallery.Files)

		// set the id for the file
		if len(gallery.Files) == 0 {
			// if this is the first file
			file.ID = 1
		} else {
			// use the next sequential number
			lid := gallery.Files[len(gallery.Files)-1]
			file.ID = lid.ID + 1
		}

		gallery.Files = append(gallery.Files, file)

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
