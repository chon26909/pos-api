// app/model/session.go
package model

import "time"

type TableSession struct {
	ID       uint64    `gorm:"primaryKey"`
	TableID  uint64    `gorm:"index;not null"`
	Token    string    `gorm:"size:36;uniqueIndex;not null"`
	IsActive bool      `gorm:"not null;default:true"`
	OpenedAt time.Time `gorm:"autoCreateTime"`
	ClosedAt *time.Time
	Table    Table `gorm:"foreignKey:TableID"`
}
