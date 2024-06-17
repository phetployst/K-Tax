package handlers

import (
	"testing"

	"github.com/KKGo-Software-engineering/assessment-tax/models"
	"github.com/stretchr/testify/assert"
)

func TestCalculateTaxAmount(t *testing.T) {
	t.Run("Taxable income less than 150000 should return tax = 0", func(t *testing.T) {
		tax, taxRefund := calculateTaxAmount(140000, 0, []models.Allowance{})
		assert.Equal(t, 0.0, tax)
		assert.Equal(t, 0.0, taxRefund)
	})

	t.Run("Taxable income between 150000 and 500000", func(t *testing.T) {
		tax, taxRefund := calculateTaxAmount(300000, 0, []models.Allowance{})
		assert.Equal(t, 15000.0, tax)
		assert.Equal(t, 0.0, taxRefund)
	})

	t.Run("Taxable income between 500000 and 1000000", func(t *testing.T) {
		tax, taxRefund := calculateTaxAmount(800000, 0, []models.Allowance{})
		assert.Equal(t, 72500.0, tax)
		assert.Equal(t, 0.0, taxRefund)
	})

	t.Run("Taxable income between 1000000 and 2000000", func(t *testing.T) {
		tax, taxRefund := calculateTaxAmount(1500000, 0, []models.Allowance{})
		assert.Equal(t, 210000.0, tax)
		assert.Equal(t, 0.0, taxRefund)
	})

	t.Run("Taxable income above 2000000", func(t *testing.T) {
		tax, taxRefund := calculateTaxAmount(2500000, 0, []models.Allowance{})
		assert.Equal(t, 485000.0, tax)
		assert.Equal(t, 0.0, taxRefund)
	})

	t.Run("Withholding tax reduces final tax", func(t *testing.T) {
		tax, taxRefund := calculateTaxAmount(2500000, 50000, []models.Allowance{})
		assert.Equal(t, 435000.0, tax)
		assert.Equal(t, 0.0, taxRefund)
	})

	t.Run("Withholding tax results in tax refund", func(t *testing.T) {
		tax, taxRefund := calculateTaxAmount(140000, 50000, []models.Allowance{})
		assert.Equal(t, 0.0, tax)
		assert.Equal(t, 50000.0, taxRefund)
	})

	t.Run("Allowances are correctly applied", func(t *testing.T) {
		allowances := []models.Allowance{
			{AllowanceType: "donation", Amount: 200000},
			{AllowanceType: "k-receipt", Amount: 70000},
			{AllowanceType: "personal", Amount: 5000},
		}
		tax, taxRefund := calculateTaxAmount(800000, 0, allowances)
		assert.Equal(t, 70000.0, tax)
		assert.Equal(t, 0.0, taxRefund)
	})
}
