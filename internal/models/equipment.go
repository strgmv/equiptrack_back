package models

import (
	"time"

	"github.com/google/uuid"
)

type Equipment struct {
	EquipmentID      uuid.UUID `json:"equipment_id" db:"equipment_id" validate:"omitempty"`
	Name             string    `json:"name,omitempty" db:"name" validate:"omitempty,lte=100"`
	ShortDescription string    `json:"short_description" db:"short_description" validate:"required,lte=200"`
	FullDescription  string    `json:"full_description" db:"full_description"`
	Reserved         bool      `json:"reserved" db:"reserved"`
	// add type
}

type EquipmentList struct {
	TotalCount int         `json:"total_count"`
	TotalPages int         `json:"total_pages"`
	Page       int         `json:"page"`
	Size       int         `json:"size"`
	HasMore    bool        `json:"has_more"`
	Equipments []Equipment `json:"equipments"`
}

type EquipmentWithUsers struct {
	EquipmentID      uuid.UUID `json:"equipment_id" db:"equipment_id" validate:"omitempty"`
	UserID           uuid.UUID `json:"user_id" db:"user_id" validate:"omitempty"`
	Name             string    `json:"name,omitempty" db:"name" validate:"omitempty,lte=100"`
	ShortDescription string    `json:"short_description" db:"short_description" validate:"required,lte=200"`
	ReservationStart time.Time `json:"reservation_start" db:"reservation_start"`
	ReservationEnd   time.Time `json:"reservation_end" db:"reservation_end"`
}

type EquipmentWithUsersList struct {
	TotalCount int                  `json:"total_count"`
	TotalPages int                  `json:"total_pages"`
	Page       int                  `json:"page"`
	Size       int                  `json:"size"`
	HasMore    bool                 `json:"has_more"`
	Equipments []EquipmentWithUsers `json:"equipments"`
}
