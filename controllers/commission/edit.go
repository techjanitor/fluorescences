package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	m "fluorescences/models"
	u "fluorescences/utils"
)

// EditController is the edit page for commission info
func EditController(c *gin.Context) {

	// holds our page metadata from settings
	metadata, err := u.GetMetadata()
	if err != nil {
		c.Error(err).SetMeta("blog.EditController.GetMetadata")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	var com m.CommissionType

	u.Storm.Get("data", "commission", &com)

	vals := struct {
		Meta       m.Metadata
		Csrf       string
		Commission m.CommissionType
	}{
		Meta:       metadata,
		Csrf:       c.MustGet("csrf_token").(string),
		Commission: com,
	}

	c.HTML(http.StatusOK, "commissionedit.tmpl", vals)

	return

}
