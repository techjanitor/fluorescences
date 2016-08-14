package controllers

import (
	"net/http"

	"github.com/asdine/storm"
	"github.com/gin-gonic/gin"

	m "fluorescences/models"
	u "fluorescences/utils"
)

// IndexController handles the categories index page
func IndexController(c *gin.Context) {
	var err error

	// holds out page metadata from settings
	metadata, err := u.GetMetadata()
	if err != nil {
		c.Error(err).SetMeta("category.IndexController.GetMetadata")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	var cats m.Categories

	// get all the categories with a limit
	err = u.Storm.All(&cats)
	if err != nil {
		c.Error(err).SetMeta("category.IndexController.All")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	for idx, cat := range cats {
		var gallery []m.GalleryType

		err = u.Storm.Find("Category", cat.ID, &gallery, storm.Limit(1))
		// if there are no galleries in the category just skip it
		if err != nil && err == storm.ErrNotFound {
			continue
		}
		if err != nil {
			c.Error(err).SetMeta("admin.PanelController.Find")
			c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
			return
		}

		cats[idx].Cover = gallery[0].Files[0].Filename
	}

	// values for template
	vals := struct {
		Meta       m.Metadata
		Categories m.Categories
	}{
		Meta:       metadata,
		Categories: cats,
	}

	c.HTML(http.StatusOK, "categories.tmpl", vals)

	return

}
