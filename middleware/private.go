package middleware

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/asdine/storm"
	"github.com/gin-gonic/gin"

	m "fluorescences/models"
	u "fluorescences/utils"
)

// Private will check the gallery status and ask for a password if its private
func Private() gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error

		// the gallery password should be given in the link
		password := c.Query("key")

		comicID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.Error(err).SetMeta("middleware.Private")
			c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
			c.Abort()
			return
		}

		var gallery m.GalleryType

		// get gallery info
		err = u.Storm.One("ID", comicID, &gallery)
		if err != nil && err != storm.ErrNotFound {
			c.Error(err).SetMeta("middleware.Private.One")
			c.HTML(http.StatusInternalServerError, "error.tmpl", nil)
			c.Abort()
			return
		}

		// handle authentication if the gallery is marked as private
		if gallery.Private {
			var authed bool

			// ask to input the password if not given
			if len(password) == 0 {
				c.Redirect(http.StatusFound, fmt.Sprintf("/gallery/key/%d", comicID))
				c.Abort()
				return
			}

			// check the given password against the gallery keys
			for _, key := range gallery.Keys {
				if key.Key == password {
					authed = true
					break
				}
			}

			if authed {
				c.Next()
			} else {
				c.Redirect(http.StatusFound, fmt.Sprintf("/gallery/key/%d", comicID))
				c.Abort()
				return
			}

		}

		c.Next()

	}
}
