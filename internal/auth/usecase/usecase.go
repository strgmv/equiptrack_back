package usecase

import (
	"context"
	"equiptrack/config"
	"equiptrack/internal/auth"
	httpErrors "equiptrack/internal/httpErrors"
	"equiptrack/internal/models"
	"equiptrack/internal/utils"
	"net/http"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type authUC struct {
	cfg      *config.Config
	authRepo auth.Repository
	logger   *logrus.Logger
}

func NewAuthUseCase(cfg *config.Config, authRepo auth.Repository, log *logrus.Logger) auth.UseCase {
	return &authUC{cfg: cfg, authRepo: authRepo, logger: log}
}

func (u *authUC) Register(ctx context.Context, user *models.User) (*models.User, error) {
	existsUser, err := u.authRepo.FindByLogin(ctx, user)
	if existsUser != nil || err == nil {
		return nil, httpErrors.NewRestErrorWithMessage(http.StatusBadRequest, httpErrors.MsgUserAlreadyExists, nil)
	}

	if err = user.PrepareCreate(); err != nil {
		return nil, httpErrors.NewBadRequestError(errors.Wrap(err, "authUC.Register.PrepareCreate"))
	}

	createdUser, err := u.authRepo.Register(ctx, user)
	if err != nil {
		return nil, err
	}
	createdUser.SanitizePassword()

	return createdUser, nil
}

func (u *authUC) Delete(ctx context.Context, userID uuid.UUID) error {
	if err := u.authRepo.Delete(ctx, userID); err != nil {
		return err
	}
	return nil
}

func (u *authUC) GetByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	user, err := u.authRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	user.SanitizePassword()

	return user, nil
}

func (u *authUC) GetUsers(ctx context.Context, pq *utils.PaginationQuery) (*models.UserList, error) {
	return u.authRepo.GetUsers(ctx, pq)
}

func (u *authUC) Login(ctx context.Context, user *models.User) (*models.UserWithToken, error) {
	foundUser, err := u.authRepo.FindByLogin(ctx, user)
	if err != nil {
		return nil, err
	}
	if err = foundUser.ComparePasswords(user.Password); err != nil {
		return nil, httpErrors.NewUnauthorizedError(errors.Wrap(err, "authUC.Login.ComparePasswords"))
	}

	foundUser.SanitizePassword()

	accessToken, err := utils.GenerateJWTToken(foundUser, u.cfg)
	if err != nil {
		return nil, httpErrors.NewInternalServerError(errors.Wrap(err, "authUC.Login.GenerateJWTToken"))
	}

	refreshToken, err := utils.NewRefreshToken()
	if err != nil {
		return nil, httpErrors.NewInternalServerError(errors.Wrap(err, "authUC.Login.NewRefreshToken"))
	}
	userWithToken := &models.UserWithToken{
		User:         foundUser,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	err = u.authRepo.SetSession(ctx, userWithToken)
	if err != nil {
		return nil, httpErrors.NewInternalServerError(errors.Wrap(err, "authUC.Login.SetSession"))
	}

	return userWithToken, nil
}

func (u *authUC) RefreshSession(ctx context.Context, userID uuid.UUID, refreshToken string) (*models.UserWithToken, error) {
	foundSession, err := u.authRepo.GetSession(ctx, userID, refreshToken)
	if err != nil {
		return nil, err
	}

	user, err := u.authRepo.GetByID(ctx, foundSession.UserID)
	if err != nil {
		return nil, err
	}

	newAccessToken, err := utils.GenerateJWTToken(user, u.cfg)
	if err != nil {
		return nil, httpErrors.NewInternalServerError(errors.Wrap(err, "authUC.RefreshSession.GenerateJWTToken"))
	}

	newRefreshToken, err := utils.NewRefreshToken()
	if err != nil {
		return nil, httpErrors.NewInternalServerError(errors.Wrap(err, "authUC.RefreshSession.NewRefreshToken"))
	}

	user.SanitizePassword()

	userWithToken := &models.UserWithToken{
		User:         user,
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}

	err = u.authRepo.SetSession(ctx, userWithToken)
	if err != nil {
		return nil, httpErrors.NewInternalServerError(errors.Wrap(err, "authUC.RefreshSession.SetSession"))
	}

	err = u.authRepo.DeleteSession(ctx, foundSession.SessionID)
	if err != nil {
		return nil, httpErrors.NewInternalServerError(errors.Wrap(err, "authUC.RefreshSession.DeleteSession"))
	}

	return userWithToken, nil
}

func (u *authUC) Logout(ctx context.Context, userID uuid.UUID, refreshToken string) error {
	err := u.authRepo.DeleteSessionByToken(ctx, userID, refreshToken)
	if err != nil {
		return httpErrors.NewInternalServerError(errors.Wrap(err, "authUC.RefreshSession.DeleteSession"))
	}
	return nil
}
