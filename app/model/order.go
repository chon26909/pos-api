// app/model/order.go
package model

import "time"

type Order struct {
	ID             uint64  `gorm:"primaryKey"`
	TableSessionID uint64  `gorm:"index;not null"`
	Status         string  `gorm:"type:enum('preparing','ready','served','checkout_requested','paid','cancelled');default:'preparing'"`
	TotalPrice     float64 `gorm:"type:decimal(10,2);not null;default:0"`
	Note           *string
	Items          []OrderItem `gorm:"foreignKey:OrderID"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type OrderItem struct {
	ID        uint64  `gorm:"primaryKey"`
	OrderID   uint64  `gorm:"index;not null"`
	ProductID uint64  `gorm:"index;not null"`
	Qty       int     `gorm:"not null;default:1"`
	UnitPrice float64 `gorm:"type:decimal(10,2);not null"`
	Status    string  `gorm:"type:enum('preparing','ready','served','cancelled');default:'preparing'"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Product Product `gorm:"foreignKey:ProductID"`
}
