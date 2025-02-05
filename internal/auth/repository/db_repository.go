package repository

import (
	"context"
	"database/sql"
	"equiptrack/internal/auth"
	"equiptrack/internal/models"
	"equiptrack/internal/utils"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type authRepo struct {
	db *sql.DB
}

// Auth Repository constructor
func NewAuthRepository(db *sql.DB) auth.Repository {
	return &authRepo{db: db}
}

func (r *authRepo) Register(ctx context.Context, user *models.User) (*models.User, error) {
	u := &models.User{}
	if err := r.db.QueryRowContext(ctx, createUserQuery, &user.Login, &user.Password, &user.Role).Scan(&u.UserID); err != nil {
		return nil, errors.Wrap(err, "authRepo.Register.StructScan")
	}

	return u, nil
}

func (r *authRepo) Delete(ctx context.Context, userID uuid.UUID) error {
	result, err := r.db.ExecContext(ctx, deleteUserQuery, userID)
	if err != nil {
		return errors.WithMessage(err, "authRepo.Delete.ExecContext")
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "authRepo.Delete.RowsAffected")
	}
	if rowsAffected == 0 {
		return errors.Wrap(sql.ErrNoRows, "authRepo.Delete.rowsAffected")
	}

	return nil
}

func (r *authRepo) GetByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	user := &models.User{}
	if err := r.db.QueryRowContext(ctx, getUserQuery, userID).Scan(
		&user.UserID,
		&user.Login,
		&user.Password,
		&user.Role,
	); err != nil {
		return nil, errors.Wrap(err, "authRepo.GetByID.QueryRowContext")
	}
	return user, nil
}

func (r *authRepo) getTotalCount(ctx context.Context) (int, error) {
	var totalCount int
	if err := r.db.QueryRowContext(ctx, qGetTotal).Scan(&totalCount); err != nil {
		return 0, errors.Wrap(err, "authRepo.getTotalCount.QueryRowContext")
	}
	return totalCount, nil
}

func (r *authRepo) GetUsers(ctx context.Context, pq *utils.PaginationQuery) (*models.UserList, error) {
	totalCount, err := r.getTotalCount(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "authRepo.GetUsers.totalCount")
	}

	if totalCount == 0 {
		return &models.UserList{
			TotalCount: totalCount,
			TotalPages: utils.GetTotalPages(totalCount, pq.GetSize()),
			Page:       pq.GetPage(),
			Size:       pq.GetSize(),
			HasMore:    utils.GetHasMore(pq.GetPage(), totalCount, pq.GetSize()),
			Users:      make([]models.User, 0),
		}, nil
	}

	rows, err := r.db.QueryContext(
		ctx,
		qGetUsers,
		pq.GetOffset(),
		pq.GetLimit(),
	)
	if err != nil {
		return nil, errors.Wrap(err, "authRepo.GetUsers.QueryContext")
	}

	var users = make([]models.User, 0, pq.GetSize())
	for rows.Next() {
		var r models.User
		err := rows.Scan(&r.UserID, &r.Login, &r.Role)
		if err != nil {
			return nil, errors.Wrap(err, "authRepo.GetUsers.QueryContext.ScanRows")
		}
		users = append(users, r)
	}

	return &models.UserList{
		TotalCount: totalCount,
		TotalPages: utils.GetTotalPages(totalCount, pq.GetSize()),
		Page:       pq.GetPage(),
		Size:       pq.GetSize(),
		HasMore:    utils.GetHasMore(pq.GetPage(), totalCount, pq.GetSize()),
		Users:      users,
	}, nil
}

func (r *authRepo) FindByLogin(ctx context.Context, user *models.User) (*models.User, error) {
	foundUser := &models.User{}
	if err := r.db.QueryRowContext(ctx, findUserByLogin, user.Login).Scan(
		&foundUser.UserID,
		&foundUser.Login,
		&foundUser.Password,
		&foundUser.Role,
	); err != nil {
		return nil, errors.Wrap(err, "authRepo.FindByLogin.QueryRowContext")
	}
	return foundUser, nil
}

func (r *authRepo) SetSession(ctx context.Context, userWithToken *models.UserWithToken) error {
	if _, err := r.db.ExecContext(ctx, setUserSession, &userWithToken.User.UserID, &userWithToken.RefreshToken); err != nil {
		return errors.Wrap(err, "authRepo.SetSession.ExecContext")
	}
	return nil
}

func (r *authRepo) GetSession(ctx context.Context, userID uuid.UUID, token string) (*models.Session, error) {
	foundSession := &models.Session{}
	if err := r.db.QueryRowContext(ctx, getUserSession, userID, token).Scan(
		&foundSession.SessionID,
		&foundSession.UserID,
		&foundSession.RefreshToken,
	); err != nil {
		return nil, errors.New("authRepo.GetSession.QueryRowContext: no token found")
	}
	return foundSession, nil
}

func (r *authRepo) DeleteSession(ctx context.Context, sessionID int) error {
	result, err := r.db.ExecContext(ctx, deleteUserSession, sessionID)
	if err != nil {
		return errors.WithMessage(err, "authRepo Delete ExecContext")
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "authRepo.Delete.RowsAffected")
	}
	if rowsAffected == 0 {
		return errors.New("authRepo.Delete.rowsAffected: no token found")
	}

	return nil
}

func (r *authRepo) DeleteSessionByToken(ctx context.Context, userID uuid.UUID, token string) error {
	result, err := r.db.ExecContext(ctx, deleteSessionByToken, userID, token)
	if err != nil {
		return errors.Wrap(err, "authRepo Delete ExecContext")
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "authRepo.Delete.RowsAffected")
	}
	if rowsAffected == 0 {
		return errors.New("authRepo.Delete.rowsAffected: no token found")
	}

	return nil
}
