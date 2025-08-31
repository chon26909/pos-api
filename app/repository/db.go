package repository

import (
	"os"
	"pos-api/app/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewDB() (*gorm.DB, error) {
	// ตัวอย่าง DSN: user:pass@tcp(127.0.0.1:3306)/orderman?charset=utf8mb4&parseTime=True&loc=Local
	dsn := os.Getenv("MYSQL_DSN")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(
		&model.Product{},
		&model.Table{},
		&model.TableSession{},
		&model.Order{},
		&model.OrderItem{},
		&model.AdminAlert{},
	); err != nil {
		return nil, err
	}

	return db, nil
}
