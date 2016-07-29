package controllers

import (
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"

	m "fluorescences/models"
	u "fluorescences/utils"
)

// ViewController handles the commission page
func ViewController(c *gin.Context) {
	var err error

	// holds out page metadata from settings
	metadata, err := u.GetMetadata()
	if err != nil {
		c.Error(err).SetMeta("image.ViewController.GetMetadata")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	var com *m.CommissionType

	err = u.Storm.Get("data", "commission", &com)
	if err != nil {
		c.Error(err).SetMeta("image.ViewController.Get")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	// make the post formatted with markdown
	unsafe := blackfriday.MarkdownCommon([]byte(com.Content))
	// sanitize the input
	html := bluemonday.UGCPolicy().SanitizeBytes(unsafe)
	// convert to template format
	com.ContentOut = template.HTML(html)
	// convert time
	com.HumanTime = com.UpdatedTime.Format("2006-01-02")

	// values for template
	vals := struct {
		Meta       m.Metadata
		Commission *m.CommissionType
	}{
		Meta:       metadata,
		Commission: com,
	}

	c.HTML(http.StatusOK, "commission.tmpl", vals)

	return

}
