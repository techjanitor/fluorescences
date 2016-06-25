package main

import (
	"flag"
	"fmt"
	"html/template"
	"net/http"
	"os"

	"github.com/boltdb/bolt"
	"github.com/facebookgo/grace/gracehttp"
	"github.com/gin-gonic/gin"

	c "fluorescences/controllers"
	u "fluorescences/utils"

	"github.com/eirka/eirka-libs/csrf"
)

func main() {

	initialize := flag.Bool("init", false, "Initialize a new database")

	flag.Parse()

	if *initialize {
		Initialize()
		os.Exit(0)
		return
	}

	t := template.Must(template.New("public").ParseGlob("templates/*.tmpl"))
	t = template.Must(t.New("admin").ParseGlob("templates/**/*.tmpl"))

	r := gin.Default()

	// load template into gin
	r.SetHTMLTemplate(t)
	r.Static("/css", "./css")
	r.Static("/images", "./images")

	// if nothing matches
	r.NoRoute(c.ErrorController)

	public := r.Group("/")
	// generates our csrf cookie
	public.Use(csrf.Cookie())

	public.GET("/", c.BlogController)
	public.GET("/blog/:page", c.BlogController)
	public.GET("/comics/:page", c.GalleryController)
	public.GET("/comic/:id/:page", c.ComicController)
	public.GET("/image/:id/:page", c.ImageController)

	admin := r.Group("/admin")

	admin.GET("/panel", c.AdminPanelController)

	admin.GET("/blog", c.BlogEditController)
	admin.POST("/blog", c.BlogPostController)

	admin.GET("/gallery", c.GalleryEditController)
	admin.POST("/gallery", c.GalleryPostController)

	s := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", "127.0.0.1", 5000),
		Handler: r,
	}

	gracehttp.Serve(s)

}

// Initialize will create a new default database
func Initialize() {
	var err error

	err = u.Bolt.Update(func(tx *bolt.Tx) (err error) {
		_, err = tx.CreateBucketIfNotExists([]byte(c.GalleryDB))
		if err != nil {
			return
		}

		_, err = tx.CreateBucketIfNotExists([]byte(c.BlogDB))
		if err != nil {
			return
		}

		_, err = tx.CreateBucketIfNotExists([]byte(u.SettingsDB))
		if err != nil {
			return
		}

		return
	})
	if err != nil {
		panic("could not init buckets")
	}

}
