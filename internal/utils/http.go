package utils

import (
	"context"
	httpErrors "equiptrack/internal/httpErrors"
	"equiptrack/internal/models"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type UserCtxKey struct{}

func GetUserFromCtx(ctx context.Context) (*models.User, error) {
	user, ok := ctx.Value(UserCtxKey{}).(*models.User)
	if !ok {
		return nil, httpErrors.Unauthorized
	}

	return user, nil
}

func GetRequestID(c echo.Context) string {
	return c.Response().Header().Get(echo.HeaderXRequestID)
}

func GetIPAddress(c echo.Context) string {
	return c.Request().RemoteAddr
}

func LogResponseError(ctx echo.Context, logger *logrus.Logger, err error) {
	logger.Errorf(
		"ErrResponseWithLog, RequestID: %s, IPAddress: %s, Error: %s",
		GetRequestID(ctx),
		GetIPAddress(ctx),
		err,
	)
}

func ErrResponseWithLog(ctx echo.Context, logger *logrus.Logger, err error) error {
	logger.Errorf(
		"ErrResponseWithLog, RequestID: %s, IPAddress: %s, Error: %s",
		GetRequestID(ctx),
		GetIPAddress(ctx),
		err,
	)
	return ctx.JSON(httpErrors.ErrorResponse(err))
}

func ReadRequest(ctx echo.Context, request interface{}) error {
	if err := ctx.Bind(request); err != nil {
		return err
	}
	return validate.StructCtx(ctx.Request().Context(), request)
}
