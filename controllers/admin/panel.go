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

	err = u.Storm.All(&galleries, storm.Reverse())
	if err != nil {
		c.Error(err).SetMeta("admin.PanelController.All")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	// get a count of the gallery images
	for idx, gallery := range galleries {
		galleries[idx].Images = len(gallery.Files)
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
	for idx, cat := range cats {
		var galleries []m.GalleryType

		err = u.Storm.Find("Category", cat.ID, &galleries)
		if err != nil && err != storm.ErrNotFound {
			c.Error(err).SetMeta("admin.PanelController.Find")
			c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
			return
		}

		cats[idx].Galleries = len(galleries)
	}

	var allblogs, blogs m.Blogs

	// get all the blogs
	err = u.Storm.All(&allblogs)
	if err != nil {
		c.Error(err).SetMeta("admin.PanelController.All")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	// we just want regular blogs in the list
	for _, blog := range allblogs {
		if blog.Notificiation {
			continue
		}
		blogs = append(blogs, blog)
	}

	var links m.Links

	// get all the links
	err = u.Storm.All(&links)
	if err != nil {
		c.Error(err).SetMeta("admin.PanelController.All")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	// values for template
	vals := struct {
		Meta       m.Metadata
		Csrf       string
		Categories m.Categories
		Galleries  m.Galleries
		Commission m.CommissionType
		Blogs      m.Blogs
		Links      m.Links
	}{
		Meta:       metadata,
		Csrf:       c.MustGet("csrf_token").(string),
		Categories: cats,
		Galleries:  galleries,
		Commission: com,
		Blogs:      blogs,
		Links:      links,
	}

	c.HTML(http.StatusOK, "panel.tmpl", vals)

	return

}
