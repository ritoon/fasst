package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"formation/model"
)

func main() {
	//os.MkdirAll("uploads", 0750)

	// sdkRickAndMorty := rickandmorty.New("https://rickandmortyapi.com")
	// caracters, err := sdkRickAndMorty.GetCaracters()
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(caracters)

	u := model.NewUser("Martin", "Solveg")
	fmt.Println(u)
	u.UpdateName("Jhon", "Lenon")
	fmt.Println(u)
	// str := model.NewUser("hello", "toto")
	// fmt.Println(str)

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Echo!")
	})

	e.POST("/users", func(c echo.Context) error {
		var u User
		if err := c.Bind(&u); err != nil {
			return c.String(400, "bad request")
		}
		return c.JSON(201, u)
	})
	e.GET("/users/:id", func(c echo.Context) error { return c.String(200, c.Param("id")) })
	e.PUT("/users/:id", func(c echo.Context) error { return c.NoContent(204) })
	e.DELETE("/users/:id", func(c echo.Context) error { return c.NoContent(204) })

	e.GET("/show", func(c echo.Context) error {
		return c.String(200, "q="+c.QueryParam("q"))
	})

	e.Static("/static", "public")

	e.GET("/hello/:name", func(c echo.Context) error {
		log.Println("hello called")
		const tmplHello = `{{define "hello"}}Hello, {{.Name}}!{{end}}`
		return c.Render(200, tmplHello, map[string]any{"Name": c.Param("name")})
	})

	e.Use(middleware.BodyLimit("10M"))

	e.POST("/upload", func(c echo.Context) error {
		fh, err := c.FormFile("file")
		if err != nil {
			return echo.NewHTTPError(400, "paramètre 'file' manquant")
		}
		src, err := fh.Open()
		if err != nil {
			return err
		}
		defer src.Close()

		dst, err := os.Create(filepath.Join("uploads", fh.Filename))
		if err != nil {
			return err
		}
		defer dst.Close()

		if _, err := io.Copy(dst, src); err != nil {
			return err
		}

		return c.JSON(201, map[string]any{
			"filename": fh.Filename,
			"size":     fh.Size,
			"path":     "/uploads/" + fh.Filename,
		})
	})

	admin := e.Group("/admin")
	admin.Use(middleware.BasicAuth(func(u, p string, c echo.Context) (bool, error) {
		return u == "joe" && p == "secret", nil
	}))
	admin.GET("/dash", func(c echo.Context) error { return c.String(200, "ok") })

	e.POST("/login", login)

	// Unauthenticated route
	e.GET("/", accessible)

	// Restricted group
	r := e.Group("/restricted")

	r.GET("/admin", restricted)

	// Configure middleware with the custom claims type
	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(jwtCustomClaims)
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

type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

var jwtSecret = []byte("secret") // à stocker en variable d'env en prod

type jwtCustomClaims struct {
	Name  string `json:"name"`
	Admin bool   `json:"admin"`
	jwt.RegisteredClaims
}

type LoginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func login(c echo.Context) error {

	var payload LoginPayload
	c.Bind(&payload)

	fmt.Println("login called", payload)

	// Throws unauthorized error
	if payload.Email != "joe@ex" || payload.Password != "Joe" {
		return echo.ErrUnauthorized
	}

	// Set custom claims
	claims := &jwtCustomClaims{
		"Jon Snow",
		true,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": t,
	})
}

func accessible(c echo.Context) error {
	return c.String(http.StatusOK, "Accessible")
}

func restricted(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwtCustomClaims)
	name := claims.Name
	return c.String(http.StatusOK, "Welcome "+name+"!")
}
