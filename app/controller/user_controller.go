package controller

import (
	"net/http"

	"github.com/dannyexmakina/go-stencil-laravalian/app/service"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

var validate = validator.New()

type UserController struct {
	users *service.UserService
}

func NewUserController(users *service.UserService) *UserController {
	return &UserController{users: users}
}

// Index は GET /users
func (ctl *UserController) Index(c echo.Context) error {
	users, err := ctl.users.List()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, users)
}

// StoreUserRequest は Laravel の FormRequest 相当
type StoreUserRequest struct {
	Name  string `json:"name" validate:"required,max=255"`
	Email string `json:"email" validate:"required,email"`
}

// Store は POST /users
func (ctl *UserController) Store(c echo.Context) error {
	var req StoreUserRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := validate.Struct(req); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}
	u, err := ctl.users.Create(req.Name, req.Email)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, u)
}
