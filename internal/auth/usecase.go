package auth

import (
	"context"
	"equiptrack/internal/models"
	"equiptrack/internal/utils"

	"github.com/google/uuid"
)

type UseCase interface {
	Register(ctx context.Context, user *models.User) (*models.User, error)
	Login(ctx context.Context, user *models.User) (*models.UserWithToken, error)
	Delete(ctx context.Context, userID uuid.UUID) error
	GetByID(ctx context.Context, userID uuid.UUID) (*models.User, error)
	RefreshSession(ctx context.Context, userID uuid.UUID, refreshToken string) (*models.UserWithToken, error)
	Logout(ctx context.Context, userID uuid.UUID, refreshToken string) error

	GetUsers(ctx context.Context, pq *utils.PaginationQuery) (*models.UserList, error)
}
