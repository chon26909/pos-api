// app/model/table.go
package model

import "time"

type Table struct {
	ID        uint64 `gorm:"primaryKey"`
	Seat      int    `gorm:"not null;default:2"`
	Status    string `gorm:"type:enum('available','occupied','closed');default:'available'"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
