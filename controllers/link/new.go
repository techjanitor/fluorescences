package controllers

import (
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"

	m "fluorescences/models"
	u "fluorescences/utils"
)

type newForm struct {
	Title string `form:"title" binding:"required"`
	URL   string `form:"url" binding:"required"`
}

// NewController adds an external link
func NewController(c *gin.Context) {
	var err error
	var inf newForm

	err = c.Bind(&inf)
	if err != nil {
		c.Error(err).SetMeta("link.NewController.Bind")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	address, err := url.Parse(inf.URL)
	if err != nil {
		c.Error(err).SetMeta("link.NewController.Parse")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	link := m.LinkType{
		Title:   inf.Title,
		Address: address.String(),
		URL:     address,
	}

	// save gallery
	err = u.Storm.Save(&link)
	if err != nil {
		c.Error(err).SetMeta("link.NewController.Save")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	c.Redirect(http.StatusFound, c.Request.Referer())

	return

}
