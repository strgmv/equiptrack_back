package auth

import (
	"context"
	"equiptrack/internal/models"
	"equiptrack/internal/utils"

	"github.com/google/uuid"
)

type Repository interface {
	Register(ctx context.Context, user *models.User) (*models.User, error)
	Delete(ctx context.Context, userID uuid.UUID) error
	GetByID(ctx context.Context, userID uuid.UUID) (*models.User, error)
	FindByLogin(ctx context.Context, user *models.User) (*models.User, error)

	SetSession(ctx context.Context, userWithToken *models.UserWithToken) error
	GetSession(ctx context.Context, userID uuid.UUID, token string) (*models.Session, error)
	DeleteSession(ctx context.Context, sessionID int) error
	DeleteSessionByToken(ctx context.Context, userID uuid.UUID, token string) error

	GetUsers(ctx context.Context, pq *utils.PaginationQuery) (*models.UserList, error)
}
