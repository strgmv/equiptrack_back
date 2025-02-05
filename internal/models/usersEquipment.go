package models

import (
	"time"

	"github.com/google/uuid"
)

type UsersEquipment struct {
	Id               int       `json:"id" db:"id" validate:"omitempty"`
	UserID           uuid.UUID `json:"user_id" db:"user_id"`
	EquipmentID      uuid.UUID `json:"equipment_id" db:"equipment_id"`
	ReservationStart time.Time `json:"reservation_start" db:"reservation_start"`
	ReservationEnd   time.Time `json:"reservation_end" db:"reservation_end"`
}

type ReservationInfo struct {
	ReservationStart time.Time `json:"reservation_start" db:"reservation_start"`
	ReservationEnd   time.Time `json:"reservation_end" db:"reservation_end"`
}

type ReservationInfoResponse struct {
	Amount          int               `json:"amount"`
	ReservationInfo []ReservationInfo `json:"data"`
}
