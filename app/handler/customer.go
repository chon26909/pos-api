// app/handler/customer.go
package handler

import (
	"net/http"
	"pos-api/app/dto"
	"pos-api/app/model"
	"pos-api/app/service"
	"time"

	"github.com/labstack/echo/v4"

	"gorm.io/gorm"
)

type CustomerHandler struct {
	DB       *gorm.DB
	OrderSvc *service.OrderService
}

func (h *CustomerHandler) Checkin(c echo.Context) error {
	var req dto.CheckinRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "bad request"})
	}

	sess, err := h.OrderSvc.EnsureActiveSession(req.TableID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	// URL สำหรับ QR = https://your.app/?token=...
	return c.JSON(http.StatusOK, echo.Map{
		"table_id":    req.TableID,
		"table_token": sess.Token,
		"qr_url":      "https://your.app/?token=" + sess.Token,
	})
}

func (h *CustomerHandler) ListProducts(c echo.Context) error {
	var ps []model.Product
	if err := h.DB.Find(&ps).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	var out []dto.ProductOut
	for _, p := range ps {
		out = append(out, dto.ProductOut{ID: p.ID, Name: p.Name, Detail: p.Detail, Price: p.Price})
	}
	return c.JSON(http.StatusOK, out)
}

func (h *CustomerHandler) CreateOrder(c echo.Context) error {
	var req dto.CreateOrderRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "bad request"})
	}
	var ins []service.OrderItemInput
	for _, it := range req.Items {
		q := it.Qty
		if q <= 0 {
			q = 1
		}
		ins = append(ins, service.OrderItemInput{ProductID: it.ProductID, Qty: q})
	}
	order, err := h.OrderSvc.CreateOrder(req.TableToken, ins, req.Note)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, echo.Map{"order_id": order.ID, "status": order.Status})
}

func (h *CustomerHandler) History(c echo.Context) error {
	var req dto.HistoryRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "bad request"})
	}

	sess, err := h.OrderSvc.GetSessionByToken(req.TableToken)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "session not found"})
	}

	var orders []model.Order
	if err := h.DB.Preload("Items").Preload("Items.Product").
		Where("table_session_id = ?", sess.ID).
		Order("id desc").
		Find(&orders).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	var resp dto.HistoryOut
	resp.TableID = sess.TableID
	var grand float64
	for _, o := range orders {
		var items []dto.OrderItemOut
		for _, it := range o.Items {
			items = append(items, dto.OrderItemOut{
				ID: it.ID, ProductID: it.ProductID, Name: it.Product.Name,
				Qty: it.Qty, UnitPrice: it.UnitPrice, Status: it.Status,
			})
		}
		resp.Orders = append(resp.Orders, dto.OrderOut{
			ID: o.ID, Status: o.Status, TotalPrice: o.TotalPrice, Note: o.Note,
			Items: items, CreatedAt: o.CreatedAt.Format(time.RFC3339),
		})
		grand += o.TotalPrice
	}
	resp.GrandTotal = grand
	return c.JSON(http.StatusOK, resp)
}

func (h *CustomerHandler) Checkout(c echo.Context) error {
	var req dto.CheckoutRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "bad request"})
	}

	if err := h.OrderSvc.MarkCheckoutRequested(req.TableToken); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, echo.Map{"message": "checkout_requested"})
}
