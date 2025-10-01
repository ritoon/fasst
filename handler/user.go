package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"

	"formation/model"
)

func (h *Handler) CreateUser(c echo.Context) error {
	// on récupère le payload
	var u model.User
	err := c.Bind(&u)
	if err != nil {
		return c.String(400, "bad request")
	}
	err = u.ValidateForCreate()
	if err != nil {
		return c.String(400, "bad request")
	}

	ctx := context.Background()

	err = h.dbUsers.CreateUser(ctx, &u)
	if err != nil {
		log.Println(err)
		return c.String(http.StatusInternalServerError, "db err")
	}

	return c.JSON(201, u)
}

func (h *Handler) UpdateUserByID(c echo.Context) error {
	payload := make(map[string]interface{})
	// model.UserForUpdate
	if err := c.Bind(&payload); err != nil {
		c.JSON(http.StatusBadRequest, nil)
	}
	fmt.Println("try to update user with this payload:", payload)

	userID := c.Param("id")

	ctx := context.Background()
	u, err := h.dbUsers.GetUserByID(ctx, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, nil)
	}

	u.UpdateFromMap(payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, nil)
	}

	err = h.dbUsers.UpdateUser(ctx, userID, u)
	if err != nil {
		c.JSON(http.StatusInternalServerError, nil)
	}

	return c.JSON(http.StatusAccepted, u)
}

func (h *Handler) GetUserByID(c echo.Context) error {
	userID := c.Param("id")

	ctx := context.Background()
	u, err := h.dbUsers.GetUserByID(ctx, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, nil)
	}

	return c.JSON(200, u)
}

func (h *Handler) DeleteUserByID(c echo.Context) error {
	userID := c.Param("id")

	ctx := context.Background()
	err := h.dbUsers.DeleteUser(ctx, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, nil)
	}

	return c.NoContent(204)
}

type LoginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type JwtCustomClaims struct {
	Name  string `json:"name"`
	Admin bool   `json:"admin"`
	jwt.RegisteredClaims
}

func (h *Handler) Login(c echo.Context) error {

	var payload LoginPayload
	c.Bind(&payload)

	fmt.Println("login called", payload)

	// Throws unauthorized error
	if payload.Email != "joe@ex" || payload.Password != "Joe" {
		return echo.ErrUnauthorized
	}

	// Set custom claims
	claims := &JwtCustomClaims{
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
