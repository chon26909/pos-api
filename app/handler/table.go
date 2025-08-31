package handler

import (
	"net/http"
	"pos-api/app/dto"
	"pos-api/app/service"
	"strconv"

	"github.com/labstack/echo/v4"
)

type AdminTableHandler struct {
	TableSvc *service.TableService
}

// POST /api/admin/table/create
func (h *AdminTableHandler) CreateTable(c echo.Context) error {
	var in dto.CreateTableRequest
	if err := c.Bind(&in); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "bad request"})
	}
	tb, err := h.TableSvc.Create(c.Request().Context(), in.Seat, in.Status)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, tb)
}

// GET /api/admin/tables
func (h *AdminTableHandler) ListTables(c echo.Context) error {
	list, err := h.TableSvc.List(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, list)
}

// PUT /api/admin/table/:id/status
func (h *AdminTableHandler) UpdateTableStatus(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil || id == 0 {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
	}
	var in dto.UpdateTableStatusRequest
	if err := c.Bind(&in); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "bad request"})
	}
	if err := h.TableSvc.UpdateStatus(c.Request().Context(), id, in.Status); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, echo.Map{"message": "updated"})
}
