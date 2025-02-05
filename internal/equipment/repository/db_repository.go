package repository

import (
	"context"
	"database/sql"
	"equiptrack/internal/equipment"
	"equiptrack/internal/models"
	"equiptrack/internal/utils"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type equipmentRepo struct {
	db *sql.DB
}

func NewEquipmentRepository(db *sql.DB) equipment.Repository {
	return &equipmentRepo{db: db}
}

func (r *equipmentRepo) Create(ctx context.Context, equipment *models.Equipment) (*models.Equipment, error) {
	e := &models.Equipment{}
	if err := r.db.QueryRowContext(
		ctx, qCreateEquipment, &equipment.Name, &equipment.ShortDescription, &equipment.FullDescription,
	).Scan(
		&e.Name, &e.ShortDescription, &e.FullDescription, &e.EquipmentID); err != nil {
		return nil, errors.Wrap(err, "equipmentRepo.Create.StructScan")
	}
	return e, nil
}

func (r *equipmentRepo) Update(ctx context.Context, equipment *models.Equipment) error {
	result, err := r.db.ExecContext(
		ctx, qUpdateEquipment, &equipment.Name, &equipment.ShortDescription, &equipment.FullDescription, &equipment.EquipmentID,
	)
	if err != nil {
		return errors.Wrap(err, "equipmentRepo.Update.ExecContext")
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "equipmentRepo.Update.RowsAffected")
	}
	if rowsAffected == 0 {
		return errors.Wrap(sql.ErrNoRows, "equipmentRepo.Update.rowsAffected")
	}
	return nil
}

func (r *equipmentRepo) Delete(ctx context.Context, equipmentID uuid.UUID) error {
	result, err := r.db.ExecContext(ctx, qDeleteEquipment, equipmentID)
	if err != nil {
		return errors.WithMessage(err, "equipmentRepo.Delete.ExecContext")
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "equipmentRepo.Delete.RowsAffected")
	}
	if rowsAffected == 0 {
		return errors.Wrap(sql.ErrNoRows, "equipmentRepo.Delete.rowsAffected")
	}

	return nil
}

func (r *equipmentRepo) GetByID(ctx context.Context, equipmentID uuid.UUID) (*models.Equipment, error) {
	equipment := &models.Equipment{}
	if err := r.db.QueryRowContext(ctx, qGetEquipment, equipmentID).Scan(
		&equipment.EquipmentID,
		&equipment.Name,
		&equipment.ShortDescription,
		&equipment.FullDescription,
	); err != nil {
		return nil, errors.Wrap(err, "equipmentRepo.GetByID.QueryRowContext")
	}
	return equipment, nil
}

func (r *equipmentRepo) getTotalCount(ctx context.Context) (int, error) {
	var totalCount int
	if err := r.db.QueryRowContext(ctx, qGetTotal).Scan(&totalCount); err != nil {
		return 0, errors.Wrap(err, "equipmentRepo.getTotalCount.QueryRowContext")
	}
	return totalCount, nil
}

func (r *equipmentRepo) getReservedByUserCount(ctx context.Context, userId uuid.UUID) (int, error) {
	var totalCount int
	if err := r.db.QueryRowContext(ctx, qGetTotalReservedByUser, userId).Scan(&totalCount); err != nil {
		return 0, errors.Wrap(err, "equipmentRepo.getReservedByUserCount.QueryRowContext")
	}
	return totalCount, nil
}

func (r *equipmentRepo) GetEquipments(ctx context.Context, pq *utils.PaginationQuery) (*models.EquipmentList, error) {
	totalCount, err := r.getTotalCount(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "equipmentRepo.GetEquipments.totalCount")
	}

	if totalCount == 0 {
		return &models.EquipmentList{
			TotalCount: totalCount,
			TotalPages: utils.GetTotalPages(totalCount, pq.GetSize()),
			Page:       pq.GetPage(),
			Size:       pq.GetSize(),
			HasMore:    utils.GetHasMore(pq.GetPage(), totalCount, pq.GetSize()),
			Equipments: make([]models.Equipment, 0),
		}, nil
	}

	rows, err := r.db.QueryContext(
		ctx,
		qGetEquipments,
		pq.GetOffset(),
		pq.GetLimit(),
	)
	if err != nil {
		return nil, errors.Wrap(err, "equipmentRepo.GetEquipments.QueryContext")
	}

	var equipments = make([]models.Equipment, 0, pq.GetSize())
	for rows.Next() {
		var r models.Equipment
		err := rows.Scan(&r.EquipmentID, &r.Name, &r.ShortDescription, &r.Reserved)
		if err != nil {
			return nil, errors.Wrap(err, "equipmentRepo.GetEquipments.QueryContext.ScanRows")
		}
		equipments = append(equipments, r)
	}

	return &models.EquipmentList{
		TotalCount: totalCount,
		TotalPages: utils.GetTotalPages(totalCount, pq.GetSize()),
		Page:       pq.GetPage(),
		Size:       pq.GetSize(),
		HasMore:    utils.GetHasMore(pq.GetPage(), totalCount, pq.GetSize()),
		Equipments: equipments,
	}, nil
}

func (r *equipmentRepo) GetUserEquipments(ctx context.Context, pq *utils.PaginationQuery, id uuid.UUID) (*models.EquipmentList, error) {
	totalCount, err := r.getReservedByUserCount(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "equipmentRepo.GetUserEquipments.totalCount")
	}
	if totalCount == 0 {
		return &models.EquipmentList{
			TotalCount: totalCount,
			TotalPages: utils.GetTotalPages(totalCount, pq.GetSize()),
			Page:       pq.GetPage(),
			Size:       pq.GetSize(),
			HasMore:    utils.GetHasMore(pq.GetPage(), totalCount, pq.GetSize()),
			Equipments: make([]models.Equipment, 0),
		}, nil
	}

	rows, err := r.db.QueryContext(
		ctx,
		qGetUserEquipments,
		pq.GetOffset(),
		pq.GetLimit(),
		id,
	)
	if err != nil {
		return nil, errors.Wrap(err, "equipmentRepo.GetUserEquipments.QueryContext")
	}

	var equipments = make([]models.Equipment, 0, pq.GetSize())
	for rows.Next() {
		var r models.Equipment
		err := rows.Scan(&r.EquipmentID, &r.Name, &r.ShortDescription, &r.Reserved)
		if err != nil {
			return nil, errors.Wrap(err, "equipmentRepo.GetUserEquipments.QueryContext.ScanRows")
		}
		equipments = append(equipments, r)
	}

	return &models.EquipmentList{
		TotalCount: totalCount,
		TotalPages: utils.GetTotalPages(totalCount, pq.GetSize()),
		Page:       pq.GetPage(),
		Size:       pq.GetSize(),
		HasMore:    utils.GetHasMore(pq.GetPage(), totalCount, pq.GetSize()),
		Equipments: equipments,
	}, nil
}

func (r *equipmentRepo) GetReservationInfo(ctx context.Context, equipmentId uuid.UUID) (*models.ReservationInfoResponse, error) {
	rows, err := r.db.QueryContext(
		ctx,
		qGetReservationInfo,
		equipmentId,
	)

	if err != nil {
		return nil, errors.Wrap(err, "equipmentRepo.GetReservationInfo.QueryContext")
	}

	var count = 0
	var info = make([]models.ReservationInfo, 0)
	for rows.Next() {
		var r models.ReservationInfo
		err := rows.Scan(&r.ReservationEnd, &r.ReservationStart)
		if err != nil {
			return nil, errors.Wrap(err, "equipmentRepo.GetReservationInfo.QueryContext.ScanRows")
		}
		count += 1
		info = append(info, r)
	}

	return &models.ReservationInfoResponse{
		Amount:          count,
		ReservationInfo: info,
	}, nil
}

func (r *equipmentRepo) IsEquipmentReservedAt(ctx context.Context, equipmentId uuid.UUID, start time.Time, end time.Time) (bool, error) {
	var busy bool
	if err := r.db.QueryRowContext(ctx, qIsReserved, equipmentId, start, end).Scan(&busy); err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return true, errors.Wrap(err, "equipmentRepo.IsEquipmentReservedAt.QueryRowContext")
	}
	return true, nil
}

func (r *equipmentRepo) ReserveEquipment(ctx context.Context, ue *models.UsersEquipment) (bool, error) {
	var reserved, err = r.IsEquipmentReservedAt(ctx, ue.EquipmentID, ue.ReservationStart, ue.ReservationEnd)
	if err != nil {
		return false, err
	}
	if reserved {
		return false, nil
	}
	_, err = r.db.ExecContext(ctx, qReserve, ue.UserID, ue.EquipmentID, ue.ReservationStart, ue.ReservationEnd)
	if err != nil {
		return false, errors.Wrap(err, "equipmentRepo.ReserveEquipment.ExecContext")
	}
	return true, nil
}
