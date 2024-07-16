package config

import (
	"os"

	"github.com/phetployst/K-Tax/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func ConnectDB() {
	dsn, ok := os.LookupEnv("DATABASE_URL")
	if !ok {
		panic("DATABASE_URL environment variable not set")
	}

	// newLogger := logger.New(
	// 	log.New(os.Stdout, "\r\n", log.LstdFlags),
	// 	logger.Config{
	// 		SlowThreshold: time.Second,
	// 		LogLevel:      logger.Info,
	// 		Colorful:      true,
	// 	},
	// )

	DB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		// Logger: newLogger,
	})
	if err != nil {
		panic("failed to connect to database")
	}

	DB.AutoMigrate(&models.TaxCalculation{}, &models.Allowance{}, &models.AdminSetting{})

	db = DB
}

func GetDB() *gorm.DB {
	return db
}
