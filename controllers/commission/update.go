package controllers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	m "fluorescences/models"
	u "fluorescences/utils"
)

type updateForm struct {
	Open    bool   `form:"open"`
	Content string `form:"content" binding:"required"`
}

// UpdateController updates commission information
func UpdateController(c *gin.Context) {
	var err error
	var uf updateForm

	err = c.Bind(&uf)
	if err != nil {
		fmt.Println(uf)
		c.Error(err).SetMeta("commission.UpdateController.Bind")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	// start transaction
	tx, err := u.Storm.Begin(true)
	if err != nil {
		c.Error(err).SetMeta("commission.UpdateController.Begin")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}
	defer tx.Rollback()

	var com m.CommissionType
	var state string

	err = tx.Get("data", "commission", &com)
	if err != nil {
		c.Error(err).SetMeta("commission.UpdateController.Get")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	com.Open = uf.Open
	com.Content = uf.Content
	com.UpdatedTime = time.Now()

	err = tx.Set("data", "commission", &com)
	if err != nil {
		c.Error(err).SetMeta("commission.UpdateController.Set")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	if com.Open {
		state = "open"
	} else {
		state = "closed"
	}

	notification := m.BlogType{
		User:          u.MustGetUsername(),
		Notificiation: true,
		StoredTime:    time.Now(),
		Title:         "Commissions Updated",
		Content:       fmt.Sprintf("<a href=\"/commission\">%s</a>", strings.ToUpper(state)),
	}

	// save blog notification
	err = tx.Save(&notification)
	if err != nil {
		c.Error(err).SetMeta("commission.UpdateController.Save")
		c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
		return
	}

	// commit
	tx.Commit()

	c.Redirect(http.StatusFound, "/commission")

	return

}
