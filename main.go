package main

import (
	"os"

	"github.com/labstack/echo/v4"
	"github.com/phetployst/K-Tax/config"
	"github.com/phetployst/K-Tax/handlers"
	"github.com/phetployst/K-Tax/middleware"
)

func main() {
	e := echo.New()

	config.ConnectDB()

	e.POST("/tax/calculations", handlers.CalculateTax)

	admin := e.Group("/admin", middleware.BasicAuth)
	admin.POST("/deductions/personal", handlers.SetPersonalDeduction)
	admin.POST("/deductions/k-receipt", handlers.SetKReceiptDeduction)

	port, exists := os.LookupEnv("PORT")
	if !exists {
		port = "8080"
	}

	e.Logger.Fatal(e.Start(":" + port))
}
