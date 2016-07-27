package main

import (
	"flag"
	"fmt"
	"html/template"
	"net/http"
	"os"

	"github.com/facebookgo/grace/gracehttp"
	"github.com/gin-gonic/gin"

	c "fluorescences/controllers"

	u "fluorescences/utils"

	blog "fluorescences/controllers/blog"
	gallery "fluorescences/controllers/gallery"
	image "fluorescences/controllers/image"

	"github.com/eirka/eirka-libs/csrf"
)

func main() {

	// make new buckets if they dont exist
	initialize := flag.Bool("init", false, "Initialize a new database")

	flag.Parse()

	if *initialize {
		Initialize()
		os.Exit(0)
		return
	}

	// load the site templates
	t := template.Must(template.New("public").ParseGlob("templates/*.tmpl"))

	r := gin.Default()

	// load template into gin
	r.SetHTMLTemplate(t)

	// serve our static files
	r.Static("/css", "./css")
	r.Static("/images", "./images")

	// if nothing matches
	r.NoRoute(c.ErrorController)

	// routing group for public handlers
	public := r.Group("/")
	// generates our csrf cookie
	public.Use(csrf.Cookie())

	public.GET("/", blog.ViewController)
	public.GET("/blog/:page", blog.ViewController)
	public.GET("/comics/:page", gallery.IndexController)
	public.GET("/comic/:id/:page", gallery.ViewController)
	public.GET("/image/:id/:page", image.ViewController)

	// routing group for admin handlers
	admin := r.Group("/admin")

	admin.GET("/panel", c.AdminPanelController)

	admin.GET("/blog", blog.NewController)
	admin.POST("/blog/new", blog.PostController)

	admin.GET("/gallery", gallery.NewController)
	admin.GET("/gallery/edit/:id", gallery.EditController)
	admin.POST("/gallery/new", gallery.PostController)
	admin.POST("/gallery/delete", gallery.DeleteController)
	admin.POST("/gallery/update", gallery.UpdateController)
	admin.POST("/gallery/image/new", image.NewController)
	admin.POST("/gallery/image/delete", image.DeleteController)

	s := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", "0.0.0.0", 5000),
		Handler: r,
	}

	gracehttp.Serve(s)

}

// Initialize will create a new default database
func Initialize() {
	var err error

	err = u.InitMetadata()
	if err != nil {
		panic("could not init metadata")
	}
}
