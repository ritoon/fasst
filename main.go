package main

import (
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"formation/db"
	"formation/db/local"
	"formation/db/mock"
	"formation/handler"
)

func buildStore() (db.UserStore, error) {
	if os.Getenv("USE_MOCK") == "1" {
		return mock.New(), nil
	}
	return local.New("app.db")
}

func main() {

	// connexion à la base de données
	store, err := buildStore()
	if err != nil {
		panic(err)
	}
	defer store.Close()

	// création du routeur
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// initialise les handlers
	handler.New(e, store)

	// Unauthenticated route
	e.GET("/", accessible)
	// Restricted group
	r := e.Group("/restricted")
	r.GET("/admin", restricted)

	// Configure middleware with the custom claims type
	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(handler.JwtCustomClaims)
		},
		SigningKey: []byte("secret"),
	}
	r.Use(echojwt.WithConfig(config))

	e.HTTPErrorHandler = func(err error, c echo.Context) {
		code := 500
		if he, ok := err.(*echo.HTTPError); ok {
			code = he.Code
		}
		_ = c.JSON(code, map[string]any{"error": err.Error()})
	}

	e.Logger.Fatal(e.Start(":1323"))

}

func accessible(c echo.Context) error {
	return c.String(http.StatusOK, "Accessible")
}

func restricted(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*handler.JwtCustomClaims)
	name := claims.Name
	return c.String(http.StatusOK, "Welcome "+name+"!")
}
