package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"

	"github.com/facebookgo/grace/gracehttp"
	"github.com/gin-gonic/gin"
	"github.com/urfave/cli"

	c "fluorescences/controllers"
	u "fluorescences/utils"

	admin "fluorescences/controllers/admin"
	blog "fluorescences/controllers/blog"
	gallery "fluorescences/controllers/gallery"
	image "fluorescences/controllers/image"

	"github.com/eirka/eirka-libs/csrf"
)

func main() {

	app := cli.NewApp()

	app.Name = "Fluorescences"
	app.Usage = "An art gallery"
	app.Version = "RC1"
	app.Copyright = "(c) 2016 Techjanitor"

	app.Commands = []cli.Command{
		{
			Name:  "init",
			Usage: "initialize a component for the first time",
			Subcommands: []cli.Command{
				{
					Name:  "user",
					Usage: "initialize the user",
					Action: func(c *cli.Context) error {
						name := c.Args().First()
						if name == "" {
							return cli.NewExitError("username required", 1)
						}
						return u.InitUser(name)
					},
				},
				{
					Name:  "secret",
					Usage: "initialize the HMAC secret",
					Action: func(c *cli.Context) error {
						return u.InitSecret()
					},
				},
			},
		},
		{
			Name:  "start",
			Usage: "start the server",
			Action: func(c *cli.Context) error {
				server()
				return nil
			},
		},
	}

	app.Run(os.Args)

}

func server() {

	// load the site templates
	t := template.Must(template.New("public").Funcs(u.TemplateFuncs).ParseGlob("templates/*.tmpl"))

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
	authed := r.Group("/admin")

	authed.GET("/login", admin.LoginController)
	authed.GET("/panel", admin.GalleryController)

	authed.GET("/blog", blog.NewController)
	authed.POST("/blog/new", blog.PostController)
	authed.POST("/blog/delete", blog.DeleteController)

	authed.GET("/gallery", gallery.NewController)
	authed.GET("/gallery/edit/:id", gallery.EditController)
	authed.POST("/gallery/new", gallery.PostController)
	authed.POST("/gallery/delete", gallery.DeleteController)
	authed.POST("/gallery/update", gallery.UpdateController)
	authed.POST("/gallery/image/new", image.NewController)
	authed.POST("/gallery/image/delete", image.DeleteController)

	s := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", "0.0.0.0", 5000),
		Handler: r,
	}

	gracehttp.Serve(s)

}
