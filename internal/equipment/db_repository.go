package equipment

import (
	"context"
	"equiptrack/internal/models"
	"equiptrack/internal/utils"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, equipment *models.Equipment) (*models.Equipment, error)
	Update(ctx context.Context, equipment *models.Equipment) error
	Delete(ctx context.Context, equipmentID uuid.UUID) error
	GetByID(ctx context.Context, equipmentID uuid.UUID) (*models.Equipment, error)
	GetEquipments(ctx context.Context, pq *utils.PaginationQuery) (*models.EquipmentList, error)
	GetUserEquipments(ctx context.Context, pq *utils.PaginationQuery, id uuid.UUID) (*models.EquipmentList, error)
	GetReservationInfo(ctx context.Context, equipmentId uuid.UUID) (*models.ReservationInfoResponse, error)
	ReserveEquipment(ctx context.Context, reservation *models.UsersEquipment) (bool, error)
}
