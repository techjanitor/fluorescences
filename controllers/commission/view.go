package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	m "fluorescences/models"
	u "fluorescences/utils"
)

// ViewController handles the commission page
func ViewController(c *gin.Context) {
	var err error

	// holds out page metadata from settings
	metadata, err := u.GetMetadata()
	if err != nil {
		c.Error(err).SetMeta("commission.ViewController.GetMetadata")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	var com *m.CommissionType

	err = u.Storm.Get("data", "commission", &com)
	if err != nil {
		c.Error(err).SetMeta("commission.ViewController.Get")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	// format post with markdown
	com.ContentOut = u.Markdown(com.Content)
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
