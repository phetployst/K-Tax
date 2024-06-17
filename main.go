package main

import (
	"os"

	"github.com/KKGo-Software-engineering/assessment-tax/config"
	"github.com/KKGo-Software-engineering/assessment-tax/handlers"
	"github.com/KKGo-Software-engineering/assessment-tax/middleware"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	config.ConnectDB()

	e.POST("/tax/calculations", handlers.CalculateTax)

	admin := e.Group("/admin", middleware.BasicAuth)
	admin.POST("/admin/deductions/personal", handlers.SetPersonalDeduction)
	admin.POST("/admin/deductions/k-receipt", handlers.SetKReceiptDeduction)

	port, exists := os.LookupEnv("PORT")
	if !exists {
		port = "8080"
	}

	e.Logger.Fatal(e.Start(":" + port))
}
