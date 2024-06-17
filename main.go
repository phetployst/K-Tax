package main

import (
	"net/http"
	"os"

	"github.com/KKGo-Software-engineering/assessment-tax/config"
	"github.com/KKGo-Software-engineering/assessment-tax/handlers"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	config.ConnectDB()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Go Bootcamp!")
	})
	e.POST("/tax/calculations", handlers.CalculateTax)

	port, exists := os.LookupEnv("PORT")
	if !exists {
		port = "8080"
	}

	e.Logger.Fatal(e.Start(":" + port))
}
