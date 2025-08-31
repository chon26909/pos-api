// go.mod เพิ่ม
// require (
//   github.com/joho/godotenv v1.5.1
// )

// app/cmd/main.go
package main

import (
	"log"
	"os"
	"pos-api/app/handler"
	"pos-api/app/repository"
	"pos-api/app/router"
	"pos-api/app/service"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	_ = godotenv.Load()

	if os.Getenv("MYSQL_DSN") == "" {
		os.Setenv("MYSQL_DSN", "root:root@tcp(127.0.0.1:3306)/orderman?charset=utf8mb4&parseTime=True&loc=Local")
	}
	if os.Getenv("APP_PORT") == "" {
		os.Setenv("APP_PORT", "8080")
	}

	db, err := repository.NewDB()
	if err != nil {
		log.Fatal(err)
	}

	// services
	orderSvc := &service.OrderService{DB: db}

	tableRepo := repository.NewTableRepository(db)
	tableSvc := service.NewTableService(tableRepo)

	// handlers
	ch := &handler.CustomerHandler{DB: db, OrderSvc: orderSvc}
	ah := &handler.AdminHandler{DB: db, OrderSvc: orderSvc}
	ath := &handler.AdminTableHandler{TableSvc: tableSvc}

	e := echo.New()
	router.Setup(e, ch, ah, ath)

	addr := os.Getenv("APP_PORT")
	log.Println("listening on", addr)
	if err := e.Start(":" + addr); err != nil {
		log.Fatal(err)
	}
}
