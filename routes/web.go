package routes

import (
	"net/http"

	"github.com/dannyexmakina/go-stencil-laravalian/app/controller"
	"github.com/dannyexmakina/go-stencil-laravalian/app/middleware"
	"github.com/dannyexmakina/go-stencil-laravalian/app/service"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// Register は routes/web.php 相当
func Register(e *echo.Echo, db *gorm.DB) {
	e.Use(middleware.RequestLogger)

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Laravalian!")
	})

	users := controller.NewUserController(service.NewUserService(db))
	e.GET("/users", users.Index)
	e.POST("/users", users.Store)
}
