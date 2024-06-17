package config

import (
	"os"

	"github.com/KKGo-Software-engineering/assessment-tax/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func ConnectDB() {
	dsn, ok := os.LookupEnv("DATABASE_URL")
	if !ok {
		panic("DATABASE_URL environment variable not set")
	}

	DB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}

	DB.AutoMigrate(&models.TaxCalculation{}, &models.Allowance{}, &models.AdminSetting{})

	db = DB
}

func GetDB() *gorm.DB {
	return db
}
