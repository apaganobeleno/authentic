package actions

import (
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/middleware"

	"github.com/apaganobeleno/authentic"
	"github.com/apaganobeleno/authentic/example/models"

	"github.com/gobuffalo/envy"

	"github.com/gobuffalo/buffalo/middleware/csrf"
	"github.com/gobuffalo/buffalo/middleware/i18n"
	"github.com/gobuffalo/packr"
)

// ENV is used to help switch settings based on where the
// application is being run. Default is "development".
var ENV = envy.Get("GO_ENV", "development")
var app *buffalo.App
var T *i18n.Translator

// App is where all routes and middleware for buffalo
// should be defined. This is the nerve center of your
// application.
func App() *buffalo.App {
	if app == nil {
		app = buffalo.Automatic(buffalo.Options{
			Env:         ENV,
			SessionName: "_example_session",
		})
		// Automatically save the session if the underlying
		// Handler does not return an error.
		app.Use(middleware.SessionSaver)

		if ENV == "development" {
			app.Use(middleware.ParameterLogger)
		}

		if ENV != "test" {
			// Protect against CSRF attacks. https://www.owasp.org/index.php/Cross-Site_Request_Forgery_(CSRF)
			// Remove to disable this.
			app.Use(csrf.Middleware)
		}

		// Wraps each request in a transaction.
		//  c.Value("tx").(*pop.PopTransaction)
		// Remove to disable this.
		app.Use(middleware.PopTransaction(models.DB))

		// Setup and use translations:
		var err error
		if T, err = i18n.New(packr.NewBox("../locales"), "en-US"); err != nil {
			app.Stop(err)
		}

		app.Use(T.Middleware())

		app.ServeFiles("/assets", packr.NewBox("../public/assets"))

		//Sets up authentic filters and handlers
		authentic.Setup(app, AuthProvider{}, authentic.Config{
			LoginPath:       "/auth/login",
			AfterLoginPath:  "/secure/home",
			AfterLogoutPath: "/",
			PublicHandlers: []buffalo.Handler{
				HomeHandler,
			},
		})

		app.GET("/", HomeHandler)
		app.GET("/secure/home", SecureHome)
	}

	return app
}
