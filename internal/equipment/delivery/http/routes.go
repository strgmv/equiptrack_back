package http

import (
	"equiptrack/internal/equipment"
	"equiptrack/internal/middleware"

	"github.com/labstack/echo/v4"
)

func MapEquipmentRoutes(equipGroup *echo.Group, h equipment.Handlers, mw *middleware.MiddlewareManager) {
	equipGroup.Use(mw.AuthJWTMiddleware)
	equipGroup.POST("/create", h.Create(), mw.IsAdminMiddleware)
	equipGroup.POST("/reserve", h.ReserveEquipment())
	equipGroup.GET("/reservations_info/:equipment_id", h.GetReservationInfo())
	equipGroup.DELETE("/:equipment_id", h.Delete(), mw.IsAdminMiddleware)
	equipGroup.PUT("/update", h.Update(), mw.IsAdminMiddleware)
	equipGroup.GET("/:equipment_id", h.GetByID())
	equipGroup.GET("", h.GetEquipments())
}
