package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
	"github.com/golang-jwt/jwt/v4"
	"github.com/valyala/fasthttp"
)

type Claims struct {
	jwt.RegisteredClaims
}

func (app *Application) AuthRequired(c *fiber.Ctx) error {
	authHeader := c.GetReqHeaders()["Authorization"]

	// sanity check
	if authHeader == "" {
		return errors.New("no auth")
	}

	// split the header on spaces
	headerParts := strings.Split(authHeader, " ")
	if len(headerParts) != 2 {
		return errors.New("invalid auth0")
	}

	// check to see if we have the word Bearer
	if headerParts[0] != "Bearer" {
		return errors.New("invalid auth1")
	}

	token := headerParts[1]

	// declare an empty claims
	claims := &Claims{}

	// parse the token
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(app.JWTSecret), nil
	})

	if err != nil {
		if strings.HasPrefix(err.Error(), "token is expired by") {
			return app.RefreshAuthRequired(c)
		}
		return err
	}

	if claims.Issuer != app.JWTIssuer {
		return errors.New("invalid issuer")
	}

	if len(claims.Audience) > 0 && claims.Audience[0] != app.Domain+"_"+app.AppID {
		return errors.New("invalid Audience")
	}

	return c.Next()
}

func (app *Application) RefreshAuthRequired(c *fiber.Ctx) error {
	csrfHeader := c.GetReqHeaders()["X-Csrf-Token"]
	// sanity check
	if csrfHeader == "" {
		return errors.New("no csrf header")
	}

	//use refresh token
	c.Request().Header.Set("Authorization", "Bearer "+csrfHeader)

	origReqBody := c.Request().Body()
	origReqHeader := &c.Request().Header

	c.Request().ResetBody()
	c.Request().SetBody(c.Response().Body())
	var origRespHeader fasthttp.ResponseHeader
	c.Response().Header.CopyTo(&origRespHeader)

	//call user auth refresh token
	if err := proxy.Do(c, app.Usr.AuthUrl+"/admin/refreshJwtauth"); err != nil {
		return err
	}
	var jsonResp = struct {
		Error   bool       `json:"error"`
		Message string     `json:"message"`
		Data    TokenPairs `json:"data,omitempty"`
	}{}

	json.Unmarshal(c.Response().Body(), &jsonResp)
	origReqHeader.Set("Authorization", "Bearer "+jsonResp.Data.Token.PlainText)
	origReqHeader.Set("X-Csrf-Token", jsonResp.Data.RefreshToken.PlainText)

	c.Request().SetBody(origReqBody)
	c.Response().Header = origRespHeader

	return c.Next()
}
