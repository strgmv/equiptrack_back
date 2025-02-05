package http

import (
	"equiptrack/internal/auth"
	"equiptrack/internal/middleware"

	"github.com/labstack/echo/v4"
)

func MapAuthRoutes(authGroup *echo.Group, h auth.Handlers, mw *middleware.MiddlewareManager) {
	authGroup.POST("/register", h.Register())
	authGroup.POST("/login", h.Login())
	authGroup.POST("/refresh", h.RefreshJWT())
	authGroup.POST("/logout", h.Logout())
	authGroup.Use(mw.AuthJWTMiddleware)
	authGroup.GET("/:user_id", h.GetUserByID())
	authGroup.GET("/status", h.CheckAuthorized())

	authGroup.GET("/all", h.GetUsers(), mw.IsAdminMiddleware)
	authGroup.DELETE("/:user_id", h.Delete(), mw.IsAdminMiddleware)
}
