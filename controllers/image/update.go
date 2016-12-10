package controllers

import (
	"fmt"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"

	m "fluorescences/models"
	u "fluorescences/utils"
)

type updateForm struct {
	Gallery int    `form:"gallery" binding:"required"`
	Image   int    `form:"image" binding:"required"`
	Desc    string `form:"desc"`
}

// UpdateController updates image information
func UpdateController(c *gin.Context) {
	var err error
	var uf updateForm

	err = c.Bind(&uf)
	if err != nil {
		c.Error(err).SetMeta("image.UpdateController.Bind")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	// start transaction
	tx, err := u.Storm.Begin(true)
	if err != nil {
		c.Error(err).SetMeta("image.UpdateController.Begin")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}
	defer tx.Rollback()

	var gallery m.GalleryType

	err = tx.One("ID", uf.Gallery, &gallery)
	if err != nil {
		c.Error(err).SetMeta("image.UpdateController.One")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	sort.Sort(gallery.Files)

	// update the image info
	for i := len(gallery.Files) - 1; i >= 0; i-- {
		file := gallery.Files[i]

		if file.ID == uf.Image {
			file.Desc = uf.Desc
			// save the item
			gallery.Files[i] = file
		}
	}

	// save with updated info
	err = tx.Save(&gallery)
	if err != nil {
		c.Error(err).SetMeta("image.UpdateController.Save")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	// commit
	tx.Commit()

	c.Redirect(http.StatusFound, fmt.Sprintf("/admin/gallery/edit/%d", uf.Gallery))

	return

}
