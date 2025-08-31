// app/handler/admin.go
package handler

import (
	"net/http"
	"pos-api/app/dto"
	"pos-api/app/model"
	"pos-api/app/service"
	"strconv"

	"github.com/labstack/echo/v4"

	"gorm.io/gorm"
)

// DTOs for snake_case JSON responses
type OrderItemProductDTO struct {
	ID     uint64  `json:"id"`
	Name   string  `json:"name"`
	Detail string  `json:"detail"`
	Price  float64 `json:"price"`
}

type OrderItemDTO struct {
	ID      uint64              `json:"id"`
	Product OrderItemProductDTO `json:"product"`
	Qty     int                 `json:"qty"`
}

type OrderDTO struct {
	ID        uint64         `json:"id"`
	TableID   uint64         `json:"table_id"`
	Status    string         `json:"status"`
	Note      string         `json:"note"`
	Items     []OrderItemDTO `json:"items"`
	CreatedAt string         `json:"created_at"`
}

type AdminAlertDTO struct {
	ID        uint64 `json:"id"`
	Message   string `json:"message"`
	IsRead    bool   `json:"is_read"`
	CreatedAt string `json:"created_at"`
}

type ProductDTO struct {
	ID     uint64  `json:"id"`
	Name   string  `json:"name"`
	Detail string  `json:"detail"`
	Price  float64 `json:"price"`
}

type AdminHandler struct {
	DB       *gorm.DB
	OrderSvc *service.OrderService
}

func (h *AdminHandler) ListOrders(c echo.Context) error {
	var orders []model.Order
	if err := h.DB.Preload("Items").Preload("Items.Product").
		Where("status IN ('preparing','ready','served','checkout_requested')").
		Order("id desc").
		Find(&orders).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	var result []dto.OrderOut
	for _, o := range orders {
		var items []dto.OrderItemOut
		for _, item := range o.Items {
			items = append(items, dto.OrderItemOut{
				ID:        item.ID,
				ProductID: item.ProductID,
				Name:      item.Product.Name,
				Qty:       item.Qty,
				UnitPrice: item.UnitPrice,
				Status:    item.Status,
			})
		}
		result = append(result, dto.OrderOut{
			ID:         o.ID,
			Status:     o.Status,
			TotalPrice: o.TotalPrice,
			Note:       o.Note,
			Items:      items,
			CreatedAt:  o.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}
	return c.JSON(http.StatusOK, result)
}

func (h *AdminHandler) Alerts(c echo.Context) error {
	var alerts []model.AdminAlert
	if err := h.DB.Where("is_read = false").Order("id desc").Find(&alerts).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	// Map to DTOs (snake_case)
	var result []dto.AdminAlertOut
	for _, a := range alerts {
		result = append(result, dto.AdminAlertOut{
			ID:             a.ID,
			TableSessionID: a.TableSessionID,
			Type:           a.Type,
			IsRead:         a.IsRead,
			CreatedAt:      a.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}
	return c.JSON(http.StatusOK, result)
}

func (h *AdminHandler) UpdateOrderStatus(c echo.Context) error {
	idStr := c.Param("orderId")
	status := c.QueryParam("status")
	id, _ := strconv.ParseUint(idStr, 10, 64)
	if status == "success" {
		status = "served"
	}
	if err := h.OrderSvc.SetOrderStatus(uint64(id), status); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, echo.Map{"message": "updated"})
}

func (h *AdminHandler) CreateProduct(c echo.Context) error {
	var in struct {
		Name   string  `json:"name"`
		Detail string  `json:"detail"`
		Price  float64 `json:"price"`
	}
	if err := c.Bind(&in); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "bad request"})
	}
	p := model.Product{Name: in.Name, Detail: in.Detail, Price: in.Price}
	if err := h.DB.Create(&p).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	resp := dto.ProductOut{
		ID:     p.ID,
		Name:   p.Name,
		Detail: p.Detail,
		Price:  p.Price,
	}
	return c.JSON(http.StatusCreated, resp)
}

// Mock ชำระเงิน: ตัดจบ session + set paid ทั้งหมด
func (h *AdminHandler) Settle(c echo.Context) error {
	var in struct {
		TableSessionID uint64 `json:"table_session_id"`
	}
	if err := c.Bind(&in); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "bad request"})
	}
	if err := h.OrderSvc.SettleAndClose(in.TableSessionID); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, echo.Map{"message": "paid"})
}
