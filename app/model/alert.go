// app/model/alert.go
package model

import "time"

type AdminAlert struct {
	ID             uint64 `gorm:"primaryKey"`
	TableSessionID uint64 `gorm:"index;not null"`
	Type           string `gorm:"type:enum('checkout_request');not null"`
	IsRead         bool   `gorm:"not null;default:false"`
	CreatedAt      time.Time
}
