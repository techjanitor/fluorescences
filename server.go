package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/facebookgo/grace/gracehttp"
	"github.com/facebookgo/pidfile"
	"github.com/gin-gonic/gin"

	c "fluorescences/controllers"
	m "fluorescences/middleware"
	u "fluorescences/utils"

	admin "fluorescences/controllers/admin"
	blog "fluorescences/controllers/blog"
	category "fluorescences/controllers/category"
	com "fluorescences/controllers/commission"
	gallery "fluorescences/controllers/gallery"
	image "fluorescences/controllers/image"
	keys "fluorescences/controllers/keys"
	link "fluorescences/controllers/link"

	"github.com/eirka/eirka-libs/csrf"
)

// start will initialize the gin server
func start(name, address, port string) {

	// init store
	u.Initialize(name)

	// create pid file
	pidfile.SetPidfilePath(fmt.Sprintf("/run/fluorescences/%s.pid", name))

	err := pidfile.Write()
	if err != nil {
		panic("Could not write pid file")
	}

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
	public.GET("/categories", category.IndexController)
	public.GET("/comics/:id/:page", gallery.IndexController)
	public.GET("/comic/:id/:page", m.Private(), gallery.ViewController)
	public.GET("/image/:id/:page", m.Private(), image.ViewController)
	public.GET("/gallery/key/:id", keys.InputController)
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

	authed.GET("/panel", admin.PanelController)

	authed.GET("/blog", blog.NewController)
	authed.GET("/blog/edit/:id", blog.EditController)

	authed.GET("/link/edit/:id", link.EditController)

	authed.GET("/category", category.NewController)
	authed.GET("/category/edit/:id", category.EditController)

	authed.GET("/gallery", gallery.NewController)
	authed.GET("/gallery/edit/:id", gallery.EditController)

	authed.GET("/image/edit/:gallery/:image", image.EditController)

	// authenticates the CSRF session token
	authed.Use(csrf.Verify())

	authed.POST("/settings/update", admin.UpdateController)

	authed.POST("/blog/new", blog.PostController)
	authed.POST("/blog/delete", blog.DeleteController)
	authed.POST("/blog/update", blog.UpdateController)

	authed.POST("/link/new", link.NewController)
	authed.POST("/link/delete", link.DeleteController)
	authed.POST("/link/update", link.UpdateController)

	authed.POST("/category/new", category.PostController)
	authed.POST("/category/delete", category.DeleteController)
	authed.POST("/category/update", category.UpdateController)

	authed.POST("/gallery/new", gallery.PostController)
	authed.POST("/gallery/delete", gallery.DeleteController)
	authed.POST("/gallery/update", gallery.UpdateController)
	authed.POST("/gallery/private", gallery.PrivateController)
	authed.POST("/gallery/image/new", image.NewController)
	authed.POST("/gallery/image/delete", image.DeleteController)
	authed.POST("/gallery/image/update", image.UpdateController)
	authed.POST("/gallery/key/new", keys.NewController)
	authed.POST("/gallery/key/delete", keys.DeleteController)

	authed.POST("/commission/update", com.UpdateController)

	authed.POST("/password/update", admin.PasswordController)

	s := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", address, port),
		Handler: r,
	}

	gracehttp.Serve(s)

}
