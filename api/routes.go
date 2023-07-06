package api

import (
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func (app *Application) registerRoutes() {

	app.FiberApp.Use(recover.New())

	conf := cors.ConfigDefault
	conf.AllowOrigins = "*"
	conf.AllowCredentials = true
	conf.AllowMethods = "GET,POST,PUT,PATCH,DELETE,OPTIONS"
	conf.AllowHeaders = "Content-Type, Accept, Authorization, X-CSRF-Token"
	app.FiberApp.Use(cors.New(conf))

	// Welcome endpoint
	app.FiberApp.Get("/", app.Home)

	// User endpoints
	app.FiberApp.Post("/user/login", app.Authenticate)
	app.FiberApp.Get("/logout", app.Logout)

	app.FiberApp.Get("/movies", app.AllMovies)
	app.FiberApp.Get("/movies/:id", app.GetAMovie)
	app.FiberApp.Get("/genres", app.AllGenres)
	app.FiberApp.Get("/genres/:id", app.MoviesByGenre)
	app.FiberApp.Post("/graph", app.MoviesGraphQL)

	//admin user
	admin := app.FiberApp.Group("/admin", app.AuthRequired)
	admin.Get("/movies", app.AllMovies)
	admin.Get("/movies/:id", app.MovieForEdit)
	admin.Patch("/movies/:id", app.UpdateMovie)
	admin.Delete("/movies/:id", app.DeleteMovie)
	admin.Put("/movies/0", app.AddMove)
}
