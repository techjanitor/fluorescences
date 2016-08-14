package controllers

import (
	"net/http"
	"sort"
	"time"

	"github.com/gin-gonic/gin"

	m "fluorescences/models"
	u "fluorescences/utils"
)

type newForm struct {
	ID int `form:"id" binding:"required"`
}

// NewController add a key to a gallery
func NewController(c *gin.Context) {
	var err error
	var inf newForm

	err = c.Bind(&inf)
	if err != nil {
		c.Error(err).SetMeta("key.NewController.Bind")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	// add the image to the gallery file slice
	err = AddKey(inf.ID)
	if err != nil {
		c.Error(err).SetMeta("key.NewController.AddKey")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	c.Redirect(http.StatusFound, c.Request.Referer())

	return

}

// AddKey will add a key to the gallery
func AddKey(gid int) (err error) {

	// start transaction
	tx, err := u.Storm.Begin(true)
	if err != nil {
		return
	}
	defer tx.Rollback()

	var gallery m.GalleryType

	err = tx.One("ID", gid, &gallery)
	if err != nil {
		return
	}

	sort.Sort(gallery.Keys)

	key := m.KeyType{
		Key:     u.GenerateRandomPassword(20),
		Created: time.Now(),
	}

	gallery.Keys = append(gallery.Keys, key)

	err = tx.Save(&gallery)
	if err != nil {
		return
	}

	// commit
	tx.Commit()

	return

}
