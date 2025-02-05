package usecase

import (
	"context"
	"equiptrack/config"
	"equiptrack/internal/equipment"
	"equiptrack/internal/models"
	"equiptrack/internal/utils"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type equipmentUC struct {
	cfg           *config.Config
	equipmentRepo equipment.Repository
	logger        *logrus.Logger
}

func NewEquipmentUseCase(cfg *config.Config, equipmentRepo equipment.Repository, log *logrus.Logger) equipment.UseCase {
	return &equipmentUC{cfg: cfg, equipmentRepo: equipmentRepo, logger: log}
}

func (u *equipmentUC) Create(ctx context.Context, equipment *models.Equipment) (*models.Equipment, error) {
	newEquip, err := u.equipmentRepo.Create(ctx, equipment)
	if err != nil {
		return nil, err
	}

	return newEquip, nil
}

func (u *equipmentUC) Update(ctx context.Context, equipment *models.Equipment) error {
	return u.equipmentRepo.Update(ctx, equipment)
}

func (u *equipmentUC) Delete(ctx context.Context, equipmentID uuid.UUID) error {
	if err := u.equipmentRepo.Delete(ctx, equipmentID); err != nil {
		return err
	}

	return nil
}

func (u *equipmentUC) GetByID(ctx context.Context, equipmentID uuid.UUID) (*models.Equipment, error) {
	equipment, err := u.equipmentRepo.GetByID(ctx, equipmentID)
	if err != nil {
		return nil, err
	}

	return equipment, nil
}

func (u *equipmentUC) GetEquipments(ctx context.Context, pq *utils.PaginationQuery) (*models.EquipmentList, error) {
	return u.equipmentRepo.GetEquipments(ctx, pq)
}

func (u *equipmentUC) GetUserEquipments(
	ctx context.Context,
	pq *utils.PaginationQuery,
	userId uuid.UUID) (*models.EquipmentList, error) {
	return u.equipmentRepo.GetUserEquipments(ctx, pq, userId)
}

func (u *equipmentUC) GetReservationInfo(ctx context.Context, equipmentId uuid.UUID) (*models.ReservationInfoResponse, error) {
	return u.equipmentRepo.GetReservationInfo(ctx, equipmentId)
}

func (u *equipmentUC) ReserveEquipment(ctx context.Context, reservation *models.UsersEquipment) (bool, error) {
	return u.equipmentRepo.ReserveEquipment(ctx, reservation)
}
