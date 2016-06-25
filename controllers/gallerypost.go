package controllers

import (
	"encoding/json"
	"html/template"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/boltdb/bolt"
	"github.com/gin-gonic/gin"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"

	u "fluorescences/utils"
)

var (
	//GalleryDB is the bucket for comic galleries
	GalleryDB = "galleries"
)

// FileType holds an image file
type FileType struct {
	ID       int
	Filename string
}

// Files is a slice of FileTypes
type Files []FileType

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
	Desc       template.HTML
	Hidden     bool
	Private    bool
	HumanTime  string
	OpenTime   time.Time
	StoredTime time.Time
	Files      Files
	Keys       []KeyType
}

// GalleryForm is the input from the blog form
type GalleryForm struct {
	Title string `form:"title" binding:"required"`
	Desc  string `form:"desc" binding:"required"`
}

// GalleryPostController posts new blogs
func GalleryPostController(c *gin.Context) {
	var err error
	var gf GalleryForm

	err = c.Bind(&gf)
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

	var dst *os.File

	dst, err = os.Create("images/" + fileheader.Filename)
	if err != nil {
		c.Error(err).SetMeta("GalleryPostController")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	defer dst.Close()

	_, err = io.Copy(dst, upload)
	if err != nil {
		c.Error(err).SetMeta("GalleryPostController")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	file := FileType{
		ID:       1,
		Filename: fileheader.Filename,
	}

	files = append(files, file)

	// make the post formatted with markdown
	unsafe := blackfriday.MarkdownCommon([]byte(gf.Desc))
	// sanitize the input
	html := bluemonday.UGCPolicy().SanitizeBytes(unsafe)
	// convert to template format
	desc := template.HTML(html)

	blog := &GalleryType{
		User:       "test",
		StoredTime: time.Now(),
		Title:      gf.Title,
		Desc:       desc,
		Hidden:     true,
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
		bucket, err := tx.CreateBucketIfNotExists([]byte(GalleryDB))
		if err != nil {
			return
		}

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
