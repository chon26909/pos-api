// app/model/product.go
package model

import "time"

type Product struct {
	ID        uint64  `gorm:"primaryKey"`
	Name      string  `gorm:"size:255;not null"`
	Detail    string  `gorm:"type:text"`
	Price     float64 `gorm:"type:decimal(10,2);not null;default:0"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
