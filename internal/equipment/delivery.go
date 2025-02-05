package equipment

import "github.com/labstack/echo/v4"

type Handlers interface {
	Create() echo.HandlerFunc
	Update() echo.HandlerFunc
	Delete() echo.HandlerFunc
	GetByID() echo.HandlerFunc
	GetEquipments() echo.HandlerFunc
	GetReservationInfo() echo.HandlerFunc
	ReserveEquipment() echo.HandlerFunc
}
