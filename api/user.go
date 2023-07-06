package api

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
	"github.com/valyala/fasthttp"
)

type User struct {
	AuthUrl           string             `json:"-"`
	UserAuth          UserAuth           `json:"user_auth" validate:"required"`
	Profile           UserPorfile        `json:"profile"`
	ThirdPartySecrets []ThirdPartySecret `json:"third_party_secrets"`
}

type UserAuth struct {
	LoginID    string     `json:"login_id" validate:"required,min=2,max=100"`
	Password   string     `json:"password" validate:"required,min=4"`
	Scope      UserScope  `json:"scope" validate:"required"`
	TokenPairs TokenPairs `json:"tokenPairs"`
}

type UserPorfile struct {
	FisrtName string  `json:"first_name"`
	LastNmae  string  `json:"last_name"`
	Email     string  `json:"email"`
	Phone     string  `json:"phone"`
	Address   Address `json:"address"`
}

type Address struct {
	Street  string `json:"street"`
	City    string `json:"city"`
	State   string `json:"state"`
	Zipcode string `json:"zip_code"`
}

type Token struct {
	PlainText string        `json:"access_token"`
	Hash      []byte        `json:"-"`
	Expiry    time.Duration `json:"expiry_time"`
}

type TokenPairs struct {
	Token        Token `json:"token" bson:"-"`
	RefreshToken Token `json:"refresh_token" bson:"-"`
}

type UserRole struct {
	RoleNmae    string `json:"role_name" validate:"required"`
	Description string `json:"description"`
}

type ThirdPartySecret struct {
	KeyName     string `json:"key_name"`
	KeyValue    string `json:"key_value"`
	Description string `json:"description"`
}

type UserScope struct {
	Domain string   `json:"user_domain" validate:"required"`
	AppID  string   `json:"user_app_id" validate:"required"`
	Role   UserRole `json:"user_role" validate:"required"`
}

type UserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (app *Application) Authenticate(c *fiber.Ctx) error {
	var Payload UserLogin
	if err := c.BodyParser(&Payload); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	app.Usr.UserAuth.LoginID = Payload.Email
	app.Usr.UserAuth.Password = Payload.Password
	app.Usr.UserAuth.Scope.Role.RoleNmae = "web_user"

	if err := c.JSON(&app.Usr); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	c.Request().ResetBody()
	c.Request().SetBody(c.Response().Body())
	var origHeader fasthttp.ResponseHeader
	c.Response().Header.CopyTo(&origHeader)

	if err := proxy.Do(c, app.Usr.AuthUrl+"/jwtauth"); err != nil {
		return err
	}

	c.Response().Header = origHeader

	return nil
}
