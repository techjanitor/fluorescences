package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/boltdb/bolt"
	"github.com/gin-gonic/gin"

	u "fluorescences/utils"
)

// Files is a slice of FileTypes
type Files []u.FileType

func (f Files) Len() int {
	return len(f)
}

func (f Files) Less(i, j int) bool {
	return f[i].ID < f[j].ID
}

func (f Files) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}

// KeyType holds a key for a private gallery
type KeyType struct {
	Key string
}

// GalleryType holds a blog post
type GalleryType struct {
	ID         int
	User       string
	Title      string
	Cover      string
	Desc       string
	Hidden     bool
	Private    bool
	HumanTime  string
	OpenTime   time.Time
	StoredTime time.Time
	Files      Files
	Keys       []KeyType
}

// GalleryNewForm is the input from the blog form
type GalleryNewForm struct {
	Title string `form:"title" binding:"required"`
	Desc  string `form:"desc" binding:"required"`
}

// GalleryPostController posts new blogs
func GalleryPostController(c *gin.Context) {
	var err error
	var gnf GalleryNewForm

	err = c.Bind(&gnf)
	if err != nil {
		c.Error(err).SetMeta("GalleryPostController.Bind")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	// slice to hold our file information
	var files Files

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

	blog := &GalleryType{
		User:       "test",
		StoredTime: time.Now(),
		Title:      gnf.Title,
		Desc:       gnf.Desc,
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
func AddGallery(blog *GalleryType) (err error) {

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
