package controllers

import (
	"encoding/json"
	"net/http"
	"sort"

	"github.com/boltdb/bolt"
	"github.com/gin-gonic/gin"

	u "fluorescences/utils"
)

// ImageNewForm is the input from the new image form
type ImageNewForm struct {
	ID int `form:"id" binding:"required"`
}

// ImageNewController posts new blogs
func ImageNewController(c *gin.Context) {
	var err error
	var inf ImageNewForm

	err = c.Bind(&inf)
	if err != nil {
		c.Error(err).SetMeta("ImageNewController.Bind")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	// Check if theres a file
	upload, fileheader, err := c.Request.FormFile("file")
	if err != nil {
		c.Error(err).SetMeta("GalleryPostController")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	file, err := u.SaveFile(upload, fileheader.Filename)
	if err != nil {
		c.Error(err).SetMeta("GalleryPostController.SaveFile")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	err = AddImage(inf.ID, file)
	if err != nil {
		c.Error(err).SetMeta("ImageNewController.AddImage")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	c.Redirect(http.StatusFound, c.Request.Referer())

	return

}

// AddImage will add a blog post
func AddImage(gid int, file u.FileType) (err error) {

	// put the tumble in the database
	err = u.Bolt.Update(func(tx *bolt.Tx) (err error) {
		b := tx.Bucket([]byte(u.GalleryDB))

		id := u.Itob(gid)

		cb := b.Cursor()

		_, v := cb.Seek(id)

		var gallery GalleryType

		err = json.Unmarshal(v, &gallery)
		if err != nil {
			return
		}

		sort.Sort(gallery.Files)

		lid := gallery.Files[len(gallery.Files)-1]

		file.ID = lid.ID + 1

		gallery.Files = append(gallery.Files, file)

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
