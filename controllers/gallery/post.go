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

// NewForm is the input from the gallery form
type NewForm struct {
	Title string `form:"title" binding:"required"`
	Desc  string `form:"desc" binding:"required"`
}

// PostController posts new blogs
func PostController(c *gin.Context) {
	var err error
	var nf NewForm

	err = c.Bind(&nf)
	if err != nil {
		c.Error(err).SetMeta("GalleryPostController.Bind")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	// slice to hold our file information
	var files m.Files

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

	// this is always the first image in the gallery
	file.ID = 1

	files = append(files, file)

	blog := &m.GalleryType{
		User:       "test",
		StoredTime: time.Now(),
		Title:      nf.Title,
		Desc:       nf.Desc,
		Files:      files,
	}

	err = AddGallery(blog)
	if err != nil {
		c.Error(err).SetMeta("GalleryPostController.AddGallery")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	c.Redirect(http.StatusFound, "/comics/1")

	return

}

// AddGallery will add a blog post
func AddGallery(blog *m.GalleryType) (err error) {

	// put the tumble in the database
	err = u.Bolt.Update(func(tx *bolt.Tx) (err error) {
		bucket := tx.Bucket([]byte(u.GalleryDB))

		// get a sequence number
		id, _ := bucket.NextSequence()

		blog.ID = int(id)

		// encode our roomconfig
		encoded, err := json.Marshal(blog)
		if err != nil {
			return
		}

		// put the blog post
		err = bucket.Put(u.Itob(blog.ID), encoded)
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
