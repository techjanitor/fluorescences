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
	m "fluorescences/middleware"
	u "fluorescences/utils"

	admin "fluorescences/controllers/admin"
	blog "fluorescences/controllers/blog"
	com "fluorescences/controllers/commission"
	gallery "fluorescences/controllers/gallery"
	image "fluorescences/controllers/image"

	"github.com/eirka/eirka-libs/csrf"
)

func main() {

	app := cli.NewApp()

	app.Name = "Fluorescences"
	app.Usage = "An art gallery blog"
	app.Version = "RC1"
	app.Copyright = "(c) 2016 Tech Janitor"

	app.Commands = []cli.Command{
		{
			Name:  "init",
			Usage: "initialize a component for the first time",
			Subcommands: []cli.Command{
				{
					Name:  "data",
					Usage: "initialize the boilerplate data",
					Action: func(c *cli.Context) error {
						return u.InitData()
					},
				},
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
				start()
				return nil
			},
		},
	}

	app.Run(os.Args)

}

// start will initialize the gin server
func start() {

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

	public.GET("/", blog.ViewController)
	public.GET("/blog/:page", blog.ViewController)
	public.GET("/comics/:page", gallery.IndexController)
	public.GET("/comic/:id/:page", gallery.ViewController)
	public.GET("/image/:id/:page", image.ViewController)
	public.GET("/commission", com.ViewController)

	// routing group for admin handlers
	authed := r.Group("/admin")
	// add a CSRF cookie and session token to requests
	authed.Use(csrf.Cookie())

	authed.GET("/login", admin.LoginController)
	authed.POST("/login", admin.AuthController)
	authed.GET("/logout", admin.LogoutController)

	// ensure the user is authenticated
	authed.Use(m.Auth())

	authed.GET("/panel", admin.GalleryController)
	authed.GET("/blog", blog.NewController)
	authed.GET("/gallery", gallery.NewController)
	authed.GET("/gallery/edit/:id", gallery.EditController)
	authed.GET("/commission", com.EditController)

	// authenticates the CSRF session token
	authed.Use(csrf.Verify())

	authed.POST("/blog/new", blog.PostController)
	authed.POST("/blog/delete", blog.DeleteController)
	authed.POST("/gallery/new", gallery.PostController)
	authed.POST("/gallery/delete", gallery.DeleteController)
	authed.POST("/gallery/update", gallery.UpdateController)
	authed.POST("/gallery/image/new", image.NewController)
	authed.POST("/gallery/image/delete", image.DeleteController)
	authed.POST("/commission/update", com.UpdateController)

	s := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", "0.0.0.0", 5000),
		Handler: r,
	}

	gracehttp.Serve(s)

}
