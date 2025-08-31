package router

import (
	"pos-api/app/handler"

	"github.com/labstack/echo/v4"
)

func Setup(e *echo.Echo, ch *handler.CustomerHandler, ah *handler.AdminHandler, ath *handler.AdminTableHandler) {
	// public/customer
	e.POST("/api/checkin", ch.Checkin)
	e.GET("/api/products", ch.ListProducts)
	e.POST("/api/order", ch.CreateOrder)
	e.POST("/api/history", ch.History)
	e.POST("/api/checkout", ch.Checkout)

	// admin/kitchen
	e.GET("/api/admin/orders", ah.ListOrders)
	e.GET("/api/admin/alert", ah.Alerts)
	e.PUT("/api/admin/orders/:orderId", ah.UpdateOrderStatus)
	e.POST("/api/admin/product/create", ah.CreateProduct)
	e.POST("/api/admin/checkout/settle", ah.Settle)

	// admin/tables
	e.POST("/api/admin/table/create", ath.CreateTable)
	e.GET("/api/admin/tables", ath.ListTables)
	e.PUT("/api/admin/table/:id/status", ath.UpdateTableStatus)
}
