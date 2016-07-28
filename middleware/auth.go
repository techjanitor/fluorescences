package middleware

import (
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

	u "fluorescences/utils"
)

// Auth is a gin middleware that checks for session cookie
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {

		// try and get the jwt cookie from the request
		cookie, err := c.Request.Cookie(u.CookieName)
		// parse jwt token if its there
		if err != http.ErrNoCookie {
			token, err := jwt.ParseWithClaims(cookie.Value, &u.TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
				return u.ValidateToken(token)
			})
			// the client side should delete any saved JWT tokens on unauth error
			if err != nil || !token.Valid {
				// delete the cookie
				http.SetCookie(c.Writer, u.DeleteCookie())
				c.Error(err).SetMeta("middleware.Auth")
				c.HTML(http.StatusUnauthorized, "error.tmpl", nil)
				c.Abort()
				return
			}
		}

		// set user data for controllers
		c.Set("authenticated", true)

		c.Next()

	}
}
