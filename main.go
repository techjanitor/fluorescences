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
	t = template.Must(t.New("admin").ParseGlob("templates/**/*.tmpl"))

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

	public.GET("/", c.BlogController)
	public.GET("/blog/:page", c.BlogController)
	public.GET("/comics/:page", c.GalleryController)
	public.GET("/comic/:id/:page", c.ComicController)
	public.GET("/image/:id/:page", c.ImageController)

	// routing group for admin handlers
	admin := r.Group("/admin")

	admin.GET("/panel", c.AdminPanelController)

	admin.GET("/new/blog", c.BlogNewController)
	admin.POST("/blog", c.BlogPostController)

	admin.GET("/new/gallery", c.GalleryNewController)
	admin.POST("/new/gallery", c.GalleryPostController)
	admin.GET("/edit/gallery/:id", c.GalleryEditController)
	admin.POST("/edit/gallery", c.GalleryUpdateController)
	admin.POST("/new/image", c.ImageNewController)

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
