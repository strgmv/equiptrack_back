package http

import (
	"equiptrack/config"
	"equiptrack/internal/auth"
	httpErrors "equiptrack/internal/httpErrors"
	"equiptrack/internal/models"
	"equiptrack/internal/utils"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type authHandlers struct {
	cfg    *config.Config
	authUC auth.UseCase
	logger *logrus.Logger
}

// NewAuthHandlers Auth handlers constructor
func NewAuthHandlers(cfg *config.Config, authUC auth.UseCase, log *logrus.Logger) auth.Handlers {
	return &authHandlers{cfg: cfg, authUC: authUC, logger: log}
}

func (h *authHandlers) Register() echo.HandlerFunc {
	return func(c echo.Context) error {
		user := &models.User{}
		if err := c.Bind(user); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		createdUser, err := h.authUC.Register(c.Request().Context(), user)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusCreated, createdUser)
	}
}

func (h *authHandlers) Login() echo.HandlerFunc {
	type Login struct {
		Login    string `json:"login" db:"login" validate:"omitempty"`
		Password string `json:"password,omitempty" db:"password" validate:"required"`
	}
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		login := &Login{}

		if err := utils.ReadRequest(c, login); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}
		userWithToken, err := h.authUC.Login(ctx, &models.User{
			Login:    login.Login,
			Password: login.Password,
		})
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, userWithToken)
	}
}

func (h *authHandlers) GetUserByID() echo.HandlerFunc {
	return func(c echo.Context) error {
		uID, err := uuid.Parse(c.Param("user_id"))
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		user, err := h.authUC.GetByID(c.Request().Context(), uID)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, user)
	}
}

func (h *authHandlers) GetUsers() echo.HandlerFunc {
	return func(c echo.Context) error {
		paginationQuery, err := utils.GetPaginationFromCtx(c)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		userList, err := h.authUC.GetUsers(c.Request().Context(), paginationQuery)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, userList)
	}
}

func (h *authHandlers) Delete() echo.HandlerFunc {
	return func(c echo.Context) error {
		uID, err := uuid.Parse(c.Param("user_id"))
		if err != nil {
			return utils.ErrResponseWithLog(c, h.logger, err)
		}

		if err = h.authUC.Delete(c.Request().Context(), uID); err != nil {
			return utils.ErrResponseWithLog(c, h.logger, err)
		}

		return c.NoContent(http.StatusOK)
	}
}

func (h *authHandlers) RefreshJWT() echo.HandlerFunc {
	type UserWithRefreshToken struct {
		UserID uuid.UUID `json:"user_id"`
		Token  string    `json:"refresh_token"`
	}
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		userWithRefreshToken := &UserWithRefreshToken{}

		if err := utils.ReadRequest(c, userWithRefreshToken); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		userWithToken, err := h.authUC.RefreshSession(ctx,
			userWithRefreshToken.UserID,
			userWithRefreshToken.Token,
		)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, userWithToken)
	}
}

func (h *authHandlers) Logout() echo.HandlerFunc {
	type UserWithRefreshToken struct {
		UserID uuid.UUID `json:"user_id"`
		Token  string    `json:"refresh_token"`
	}
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		userWithRefreshToken := &UserWithRefreshToken{}

		if err := utils.ReadRequest(c, userWithRefreshToken); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		err := h.authUC.Logout(ctx,
			userWithRefreshToken.UserID,
			userWithRefreshToken.Token,
		)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}
		return c.NoContent(http.StatusOK)
	}
}

func (h *authHandlers) CheckAuthorized() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	}
}
