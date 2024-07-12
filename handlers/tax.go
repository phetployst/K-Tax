package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/KKGo-Software-engineering/assessment-tax/config"
	"github.com/KKGo-Software-engineering/assessment-tax/models"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func CalculateTax(c echo.Context) error {
	var taxCalc models.TaxCalculation
	var db = config.GetDB()

	if err := c.Bind(&taxCalc); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "cannot parse json"})
	}

	taxCalc.Tax, taxCalc.TaxRefund = calculateTaxAmount(taxCalc.TotalIncome, taxCalc.WHT, taxCalc.Allowances)

	if err := db.Create(&taxCalc).Error; err != nil {
		log.Printf("Error saving to database: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to save calculation"})
	}

	if taxCalc.Tax > 0 {
		return c.JSON(http.StatusOK, map[string]float64{"tax": taxCalc.Tax})
	} else {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"tax":       taxCalc.Tax,
			"taxRefund": taxCalc.TaxRefund,
		})
	}
}

func calculateTaxAmount(totalIncome float64, wht float64, allowances []models.Allowance) (float64, float64) {

	totalAllowances := calculateAllowance(allowances)
	for _, allowance := range allowances {
		totalAllowances += allowance.Amount
	}

	taxableIncome := totalIncome - totalAllowances

	tax := calculateProgressiveRate(taxableIncome)
	finalTax := tax - wht

	var taxRefund float64
	if finalTax < 0 {
		taxRefund = -finalTax
		finalTax = 0
	}

	return finalTax, taxRefund
}

func getAdminDeduction(db *gorm.DB) (models.AdminSetting, error) {
	var deduction models.AdminSetting
	err := db.Order("created_at desc").First(&deduction).Error
	if err != nil {
		return deduction, err
	}
	return deduction, nil
}

func calculateAllowance(allowances []models.Allowance) float64 {
	var db = config.GetDB()
	deductions, err := getAdminDeduction(db)
	if err != nil {
		fmt.Printf("Error fetching admin deductions: %v\n", err)
		return 0
	}

	for i, allowance := range allowances {
		switch allowance.AllowanceType {
		case "donation":
			if allowance.Amount > 100000 {
				allowances[i].Amount = 100000
			}
		case "k-receipt":
			if allowance.Amount > deductions.KReceiptDeduction {
				allowances[i].Amount = deductions.KReceiptDeduction
			} else if allowance.Amount > 50000 {
				allowances[i].Amount = 50000
			}
		case "personal":
			if allowance.Amount > deductions.PersonalDeduction {
				allowances[i].Amount = deductions.PersonalDeduction
			} else if allowance.Amount < 60000 {
				allowances[i].Amount = 0
			}
		case "":
			allowances[i].Amount = 0
		default:
			fmt.Printf("Warning: Unknown AllowanceType '%s' for allowance at index %d\n", allowance.AllowanceType, i)
		}
	}

	totalAllowances := float64(0)
	for _, allowance := range allowances {
		totalAllowances += allowance.Amount
	}

	return totalAllowances
}

func calculateProgressiveRate(taxableIncome float64) float64 {
	var tax float64
	switch {
	case taxableIncome <= 150000:
		tax = 0
	case taxableIncome <= 500000:
		tax = (taxableIncome - 150000) * 0.10
	case taxableIncome <= 1000000:
		tax = (taxableIncome-500000)*0.15 + 35000
	case taxableIncome <= 2000000:
		tax = (taxableIncome-1000000)*0.20 + 110000
	default:
		tax = (taxableIncome-2000000)*0.35 + 310000
	}
	return tax
}

var data struct {
	Amount float64 `json:"amount"`
}

func SetPersonalDeduction(c echo.Context) error {

	var setting models.AdminSetting
	var db = config.GetDB()

	if err := c.Bind(&data); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "cannot parse json"})
	}

	if data.Amount < 10000 || data.Amount > 100000 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid amount"})
	}

	db.FirstOrCreate(&setting)
	setting.PersonalDeduction = data.Amount

	if err := db.Save(&setting).Error; err != nil {
		log.Printf("Error saving to database: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to save personal deduction"})
	}

	return c.JSON(http.StatusOK, map[string]float64{"personalDeduction": setting.PersonalDeduction})
}

func SetKReceiptDeduction(c echo.Context) error {

	var setting models.AdminSetting
	var db = config.GetDB()

	if err := c.Bind(&data); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "cannot parse json"})
	}

	if data.Amount < 0 || data.Amount > 100000 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid amount"})
	}

	db.FirstOrCreate(&setting)
	setting.KReceiptDeduction = data.Amount

	if err := db.Save(&setting).Error; err != nil {
		log.Printf("Error saving to database: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to save KReceipt deduction"})
	}

	return c.JSON(http.StatusOK, map[string]float64{"KReceipt": setting.KReceiptDeduction})
}
