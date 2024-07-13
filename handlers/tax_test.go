package handlers

import (
	"testing"

	"github.com/KKGo-Software-engineering/assessment-tax/config"
	"github.com/KKGo-Software-engineering/assessment-tax/models"
	"github.com/go-playground/assert/v2"
)

func TestCalculateProgressiveRate(t *testing.T) {
	tests := []struct {
		netIncome float64
		expected  float64
	}{
		{netIncome: 0.0, expected: 0.0},
		{netIncome: 75000.0, expected: 0.0},
		{netIncome: 150000.0, expected: 0.0},
		{netIncome: 200000.0, expected: 5000.0},
		{netIncome: 300000.0, expected: 15000.0},
		{netIncome: 500000.0, expected: 35000.0},
		{netIncome: 600000.0, expected: 50000.0},
		{netIncome: 800000.0, expected: 80000.0},
		{netIncome: 1000000.0, expected: 110000.0},
		{netIncome: 1200000.0, expected: 150000.0},
		{netIncome: 1500000.0, expected: 210000.0},
		{netIncome: 2000000.0, expected: 310000.0},
		{netIncome: 2500000.0, expected: 485000.0},
	}

	for _, test := range tests {
		result := calculateProgressiveRate(test.netIncome)
		if result != test.expected {
			t.Errorf("For netIncome %.2f, tax rate expected = %.2f but got = %.2f", test.netIncome, test.expected, result)
		}
	}
}

// func setupTestDB() *gorm.DB {
// 	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
// 	if err != nil {
// 		panic(fmt.Sprintf("Failed to open database: %v", err))
// 	}
// 	db.AutoMigrate(&models.TaxCalculation{}, &models.Allowance{}, &models.AdminSetting{})

// 	// Insert admin settings
// 	adminSetting := models.AdminSetting{
// 		KReceiptDeduction: 50000.0,
// 		PersonalDeduction: 60000.0,
// 	}
// 	db.Create(&adminSetting)

// 	// config.ConnectDB(db) // Set the mock database to the global config

// 	return db
// }

func TestCalculateAllowance(t *testing.T) {

	t.Run("test maximum allowances", func(t *testing.T) {
		// I'll be back to mock DB
		config.ConnectDB()

		allowances := []models.Allowance{
			{AllowanceType: "donation", Amount: 120000.0},
			{AllowanceType: "k-receipt", Amount: 70000.0},
			{AllowanceType: "personal", Amount: 50000.0},
			{AllowanceType: "", Amount: 20000.0},
		}

		expected := 100000.0 + 50000.0 + 0.0

		got := calculateAllowance(allowances)

		assert.Equal(t, expected, got)
	})

	t.Run("test normal allowances", func(t *testing.T) {
		// I'll be back to mock DB
		config.ConnectDB()

		allowances := []models.Allowance{
			{AllowanceType: "donation", Amount: 90000.0},
			{AllowanceType: "k-receipt", Amount: 40000.0},
			{AllowanceType: "personal", Amount: 0.0},
		}

		expected := 90000.0 + 40000.0

		got := calculateAllowance(allowances)

		assert.Equal(t, expected, got)
	})

	t.Run("test normal allowances", func(t *testing.T) {
		// I'll be back to mock DB
		config.ConnectDB()

		allowances := []models.Allowance{
			{AllowanceType: "donation", Amount: 90000.0},
			{AllowanceType: "k-receipt", Amount: 40000.0},
			{AllowanceType: "personal", Amount: 70000.0},
		}

		expected := 90000.0 + 40000.0 + 70000.0

		got := calculateAllowance(allowances)

		assert.Equal(t, expected, got)
	})

	t.Run("non calculate personal allowance with less than 60000", func(t *testing.T) {
		// I'll be back to mock DB
		config.ConnectDB()

		allowances := []models.Allowance{
			{AllowanceType: "donation", Amount: 90000.0},
			{AllowanceType: "k-receipt", Amount: 40000.0},
			{AllowanceType: "personal", Amount: 50000.0},
		}

		expected := 90000.0 + 40000.0 + 0.0

		got := calculateAllowance(allowances)

		assert.Equal(t, expected, got)
	})

}
