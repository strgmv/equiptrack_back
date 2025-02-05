package middleware

import (
	httpErrors "equiptrack/internal/httpErrors"
	"equiptrack/internal/utils"

	"github.com/labstack/echo/v4"
)

func (mw *MiddlewareManager) IsAdminMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		u, err := utils.GetUserFromCtx(c.Request().Context())

		if err != nil || u.Role != "admin" {
			return utils.ErrResponseWithLog(c, mw.logger, httpErrors.Forbidden)
		}
		return next(c)
	}
}
