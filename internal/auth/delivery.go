package auth

import "github.com/labstack/echo/v4"

type Handlers interface {
	Register() echo.HandlerFunc
	Login() echo.HandlerFunc
	Delete() echo.HandlerFunc
	RefreshJWT() echo.HandlerFunc
	Logout() echo.HandlerFunc
	GetUserByID() echo.HandlerFunc
	CheckAuthorized() echo.HandlerFunc

	GetUsers() echo.HandlerFunc
}
