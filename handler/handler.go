package handler

import (
	"github.com/labstack/echo/v4"

	"formation/db"
)

type Handler struct {
	router  *echo.Echo
	dbUsers db.UserStore
}

func New(r *echo.Echo, db db.UserStore) *Handler {
	h := &Handler{
		router:  r,
		dbUsers: db,
	}
	h.initRoutes()
	return h
}

func (h *Handler) initRoutes() {
	h.router.POST("/users", h.CreateUser)
	h.router.GET("/users/:id", h.GetUserByID)
	h.router.PATCH("/users/:id", h.UpdateUserByID)
	h.router.DELETE("/users/:id", h.DeleteUserByID)
	h.router.POST("/login", h.Login)
}
