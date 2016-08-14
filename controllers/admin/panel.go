package controllers

import (
	"net/http"

	"github.com/asdine/storm"
	"github.com/gin-gonic/gin"

	m "fluorescences/models"
	u "fluorescences/utils"
)

// PanelController handles the admin panel
func PanelController(c *gin.Context) {
	var err error

	// holds out page metadata from settings
	metadata, err := u.GetMetadata()
	if err != nil {
		c.Error(err).SetMeta("admin.PanelController")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	var galleries m.Galleries

	u.Storm.All(&galleries, storm.Reverse())

	// get a count of the gallery images
	for _, gallery := range galleries {
		gallery.Images = len(gallery.Files)
	}

	var com m.CommissionType

	u.Storm.Get("data", "commission", &com)

	var cats m.Categories

	// get all the categories
	err = u.Storm.All(&cats)
	if err != nil {
		c.Error(err).SetMeta("admin.PanelController.All")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	// get a total of galleries in the category
	for _, cat := range cats {
		var galleries []m.GalleryType

		err = u.Storm.Find("Category", cat.ID, &galleries)
		if err != nil && err != storm.ErrNotFound {
			c.Error(err).SetMeta("admin.PanelController.Find")
			c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
			return
		}

		cat.Galleries = len(galleries)
	}

	// values for template
	vals := struct {
		Meta       m.Metadata
		Csrf       string
		Categories m.Categories
		Galleries  m.Galleries
		Commission m.CommissionType
	}{
		Meta:       metadata,
		Csrf:       c.MustGet("csrf_token").(string),
		Categories: cats,
		Galleries:  galleries,
		Commission: com,
	}

	c.HTML(http.StatusOK, "panel.tmpl", vals)

	return

}
