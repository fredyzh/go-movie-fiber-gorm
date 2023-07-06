package api

import (
	"fmt"
	"log"
	"movie/repositories/gormdb"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type Application struct {
	Domain    string
	AppID     string
	FiberApp  *fiber.App
	Usr       *User
	Validator *validator.Validate
	DSN       string
	DB        *gormdb.PostgresDBRepo
	JWTSecret string
	JWTIssuer string
}

func (app *Application) StartApp() {
	port := os.Getenv("WEB_PORT")

	app.Validator = validator.New()

	app.Domain = os.Getenv("DOMAIN")
	app.AppID = os.Getenv("APP_ID")
	app.DSN = os.Getenv("POSTGRES_DSN")
	app.Usr = &User{}
	app.Usr.AuthUrl = os.Getenv("Auth_USR_BACKEND")
	app.Usr.UserAuth = UserAuth{}
	app.Usr.UserAuth.Scope = UserScope{
		Domain: app.Domain,
		AppID:  app.AppID,
	}
	app.Usr.ThirdPartySecrets = []ThirdPartySecret{}
	app.Usr.ThirdPartySecrets = append(app.Usr.ThirdPartySecrets, ThirdPartySecret{
		KeyName: os.Getenv("JWT_KEY"),
	})
	app.JWTSecret = os.Getenv("JWT_SECRET")
	app.JWTIssuer = os.Getenv("JWT_ISSUER")

	// connect to the database
	conn, err := app.connectToDB()
	if err != nil {
		log.Fatal(err)
	}
	app.DB = &gormdb.PostgresDBRepo{DB: conn}

	//init fiber
	app.FiberApp = fiber.New()

	//register routers
	app.registerRoutes()
	log.Println("Starting application on port", port)
	app.FiberApp.Listen(fmt.Sprintf(":%s", port))
}
