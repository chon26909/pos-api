// app/dto/response.go
package dto

type ProductOut struct {
	ID     uint64  `json:"id"`
	Name   string  `json:"name"`
	Detail string  `json:"detail"`
	Price  float64 `json:"price"`
}

type OrderItemOut struct {
	ID        uint64  `json:"id"`
	ProductID uint64  `json:"product_id"`
	Name      string  `json:"name"`
	Qty       int     `json:"qty"`
	UnitPrice float64 `json:"unit_price"`
	Status    string  `json:"status"`
}

type OrderOut struct {
	ID         uint64         `json:"id"`
	Status     string         `json:"status"`
	TotalPrice float64        `json:"total_price"`
	Note       *string        `json:"note"`
	Items      []OrderItemOut `json:"items"`
	CreatedAt  string         `json:"created_at"`
}

type HistoryOut struct {
	TableID    uint64     `json:"table_id"`
	Orders     []OrderOut `json:"orders"`
	GrandTotal float64    `json:"grand_total"`
}

type AdminAlertOut struct {
	ID             uint64 `json:"id"`
	TableSessionID uint64 `json:"table_session_id"`
	Type           string `json:"type"`
	IsRead         bool   `json:"is_read"`
	CreatedAt      string `json:"created_at"`
}
