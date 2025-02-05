package http

import (
	"equiptrack/config"
	"equiptrack/internal/equipment"
	httpErrors "equiptrack/internal/httpErrors"
	"equiptrack/internal/models"
	"equiptrack/internal/utils"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type equipmentHandlers struct {
	cfg         *config.Config
	equipmentUC equipment.UseCase
	logger      *logrus.Logger
}

func NewEquipmentHandlers(cfg *config.Config, equipmentUC equipment.UseCase, log *logrus.Logger) equipment.Handlers {
	return &equipmentHandlers{cfg: cfg, equipmentUC: equipmentUC, logger: log}
}

func (h *equipmentHandlers) Create() echo.HandlerFunc {
	return func(c echo.Context) error {
		equipment := &models.Equipment{}
		if err := c.Bind(equipment); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		createdEquipment, err := h.equipmentUC.Create(c.Request().Context(), equipment)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}
		return c.JSON(http.StatusCreated, createdEquipment)
	}
}

func (h *equipmentHandlers) Update() echo.HandlerFunc {
	return func(c echo.Context) error {
		equipment := &models.Equipment{}
		if err := c.Bind(equipment); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		err := h.equipmentUC.Update(c.Request().Context(), equipment)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}
		return c.NoContent(http.StatusOK)
	}
}

func (h *equipmentHandlers) Delete() echo.HandlerFunc {
	return func(c echo.Context) error {
		eID, err := uuid.Parse(c.Param("equipment_id"))
		if err != nil {
			return utils.ErrResponseWithLog(c, h.logger, err)
		}

		if err = h.equipmentUC.Delete(c.Request().Context(), eID); err != nil {
			return utils.ErrResponseWithLog(c, h.logger, err)
		}

		return c.NoContent(http.StatusOK)
	}
}

func (h *equipmentHandlers) GetByID() echo.HandlerFunc {
	return func(c echo.Context) error {
		eID, err := uuid.Parse(c.Param("equipment_id"))
		if err != nil {
			return utils.ErrResponseWithLog(c, h.logger, err)
		}

		equipment, err := h.equipmentUC.GetByID(c.Request().Context(), eID)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, equipment)
	}
}

func (h *equipmentHandlers) GetEquipments() echo.HandlerFunc {
	return func(c echo.Context) error {
		// time.Sleep(5 * time.Second)
		paginationQuery, err := utils.GetPaginationFromCtx(c)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		id_str := c.QueryParam("user_id")
		if id_str == "" {
			// return all equipments
			equipmentList, err := h.equipmentUC.GetEquipments(c.Request().Context(), paginationQuery)
			if err != nil {
				utils.LogResponseError(c, h.logger, err)
				return c.JSON(httpErrors.ErrorResponse(err))
			}

			return c.JSON(http.StatusOK, equipmentList)
		}
		id, err := uuid.Parse(id_str)

		// return equip list for given user
		if err != nil {
			return utils.ErrResponseWithLog(c, h.logger, httpErrors.BadQueryParams)
		}
		equipmentList, err := h.equipmentUC.GetUserEquipments(c.Request().Context(), paginationQuery, id)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, equipmentList)
	}
}

func (h *equipmentHandlers) GetReservationInfo() echo.HandlerFunc {
	return func(c echo.Context) error {
		eID, err := uuid.Parse(c.Param("equipment_id"))
		if err != nil {
			return utils.ErrResponseWithLog(c, h.logger, err)
		}

		data, err := h.equipmentUC.GetReservationInfo(c.Request().Context(), eID)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, data)
	}
}

func (h *equipmentHandlers) ReserveEquipment() echo.HandlerFunc {
	return func(c echo.Context) error {
		usersEquipment := &models.UsersEquipment{}
		if err := c.Bind(usersEquipment); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		u, err := utils.GetUserFromCtx(c.Request().Context())
		if err != nil {
			return utils.ErrResponseWithLog(c, h.logger, httpErrors.Forbidden)
		}

		usersEquipment.UserID = u.UserID

		if usersEquipment.ReservationEnd.Before(usersEquipment.ReservationStart) {
			return c.JSON(httpErrors.ErrorResponse(httpErrors.BadRequest))
		}

		created, err := h.equipmentUC.ReserveEquipment(c.Request().Context(), usersEquipment)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}
		if !created {
			return c.JSON(httpErrors.ErrorResponse(httpErrors.BadRequest))
		}
		return c.NoContent(http.StatusCreated)
	}
}
