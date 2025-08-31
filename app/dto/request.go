// app/dto/request.go
package dto

type CheckinRequest struct {
	TableID uint64 `json:"table_id" validate:"required"`
}

type OrderItemIn struct {
	ProductID uint64 `json:"product_id" validate:"required"`
	Qty       int    `json:"qty" validate:"min=1"`
}

type CreateOrderRequest struct {
	TableToken string        `json:"table_token" validate:"required"`
	Note       *string       `json:"note"`
	Items      []OrderItemIn `json:"items" validate:"min=1,dive"`
}

type HistoryRequest struct {
	TableToken string `json:"table_token" validate:"required"`
}

type CheckoutRequest struct {
	TableToken string `json:"table_token" validate:"required"`
}
