package middleware

import (
	"context"
	"equiptrack/config"
	"equiptrack/internal/auth"
	httpErrors "equiptrack/internal/httpErrors"
	"equiptrack/internal/utils"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (mw *MiddlewareManager) AuthJWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenString := c.Request().Header.Get("Authorization")

		if tokenString != "" {
			if err := mw.validateJWTToken(tokenString, mw.authUC, c, mw.cfg); err != nil {
				return utils.ErrResponseWithLog(c, mw.logger, httpErrors.InvalidJWTToken)
			}
			return next(c)
		}
		return c.JSON(http.StatusUnauthorized, httpErrors.NewUnauthorizedError(httpErrors.Unauthorized))
	}
}

func (mw *MiddlewareManager) validateJWTToken(tokenString string, authUC auth.UseCase, c echo.Context, cfg *config.Config) error {
	if tokenString == "" {
		return httpErrors.InvalidJWTToken
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signin method %v", token.Header["alg"])
		}
		secret := []byte(cfg.Server.JwtSecretKey)
		return secret, nil
	})
	if err != nil {
		return err
	}

	if !token.Valid {
		return httpErrors.InvalidJWTToken
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		et, err := claims.GetExpirationTime()
		if err != nil || (*et).Time.Unix() < time.Now().Unix() {
			return httpErrors.InvalidJWTToken
		}

		userID, ok := claims["id"].(string)
		if !ok {
			return httpErrors.InvalidJWTClaims
		}

		userUUID, err := uuid.Parse(userID)
		if err != nil {
			return err
		}

		u, err := authUC.GetByID(c.Request().Context(), userUUID)
		if err != nil {
			return err
		}

		c.Set("user", u)

		ctx := context.WithValue(c.Request().Context(), utils.UserCtxKey{}, u)
		c.SetRequest(c.Request().WithContext(ctx))
	}
	return nil
}
