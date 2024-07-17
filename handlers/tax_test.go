package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/phetployst/K-Tax/config"
	"github.com/phetployst/K-Tax/models"
	"github.com/stretchr/testify/assert"
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

func TestCalculateAllowance(t *testing.T) {
	tests := []struct {
		name       string
		allowances []models.Allowance
		expected   float64
	}{
		{
			name: "Case 1: test maximum allowances",
			allowances: []models.Allowance{
				{AllowanceType: "donation", Amount: 120000.0},
				{AllowanceType: "k-receipt", Amount: 70000.0},
				{AllowanceType: "personal", Amount: 50000.0},
				{AllowanceType: "", Amount: 20000.0},
			},
			expected: 150000.0,
		},
		{
			name: "Case 2: test normal allowances 1",
			allowances: []models.Allowance{
				{AllowanceType: "donation", Amount: 90000.0},
				{AllowanceType: "k-receipt", Amount: 40000.0},
				{AllowanceType: "personal", Amount: 0.0},
			},
			expected: 130000.0,
		},
		{
			name: "Case 3: test normal allowances 2",
			allowances: []models.Allowance{
				{AllowanceType: "donation", Amount: 90000.0},
				{AllowanceType: "k-receipt", Amount: 40000.0},
				{AllowanceType: "personal", Amount: 70000.0},
			},
			expected: 200000.0,
		},
		{
			name: "Case 4: non calculate personal allowance with less than 60000",
			allowances: []models.Allowance{
				{AllowanceType: "donation", Amount: 90000.0},
				{AllowanceType: "k-receipt", Amount: 40000.0},
				{AllowanceType: "personal", Amount: 50000.0},
			},
			expected: 130000.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// I'll be back to mock DB
			config.ConnectDB()

			got := calculateAllowance(tt.allowances)

			assert.Equal(t, tt.expected, got, "total allowance got = %f but expected = %f", got, tt.expected)
		})
	}
}

func TestCalculateTaxAmount(t *testing.T) {

	tests := []struct {
		name           string
		totalIncome    float64
		wht            float64
		allowances     []models.Allowance
		expectedTax    float64
		expectedRefund float64
	}{
		{
			name:           "Case 1: Total Income 1000000, WHT 100000, Personal Allowance 60000",
			totalIncome:    1000000,
			wht:            0,
			allowances:     []models.Allowance{{AllowanceType: "personal", Amount: 60000}},
			expectedTax:    92000,
			expectedRefund: 0,
		},
		{
			name:           "Case 2: Total Income 2000000, WHT 300000, Donation Allowance 85000",
			totalIncome:    2000000,
			wht:            300000,
			allowances:     []models.Allowance{{AllowanceType: "donation", Amount: 85000}},
			expectedTax:    0,
			expectedRefund: 24000,
		},
		{
			name:           "Case 3: Total Income 500000, WHT 70000, K-receipt Allowance 70000",
			totalIncome:    500000,
			wht:            9000,
			allowances:     []models.Allowance{{AllowanceType: "k-receipt", Amount: 70000}},
			expectedTax:    16000,
			expectedRefund: 0,
		},
		{
			name:        "Case 4: Total Income 500000, WHT 100000, Multiple Allowances",
			totalIncome: 500000,
			wht:         100000,
			allowances: []models.Allowance{
				{AllowanceType: "donation", Amount: 60000},
				{AllowanceType: "personal", Amount: 8000},
				{AllowanceType: "k-receipt", Amount: 78000},
			},
			expectedTax:    0,
			expectedRefund: 91000,
		},
		{
			name:        "Case 5: Total Income 500000, WHT 100000, Multiple Allowances",
			totalIncome: 500000,
			wht:         100000,
			allowances: []models.Allowance{
				{AllowanceType: "donation", Amount: 20000},
				{AllowanceType: "personal", Amount: 85000},
				{AllowanceType: "k-receipt", Amount: 9000},
			},
			expectedTax:    0,
			expectedRefund: 84800,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tax, refund := calculateTaxAmount(test.totalIncome, test.wht, test.allowances)
			assert.Equal(t, test.expectedTax, tax, "Tax should be equal %f but got %f", test.expectedTax, tax)
			assert.Equal(t, test.expectedRefund, refund, "Refund should be equal %f but got %f", test.expectedRefund, refund)
		})
	}
}

func TestSetPersonalDeduction(t *testing.T) {
	e := echo.New()
	reqBody := `{"amount": 50000}`
	req := httptest.NewRequest(http.MethodPost, "/admin/deductions/personal", bytes.NewReader([]byte(reqBody)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, SetPersonalDeduction(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		var response map[string]float64
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, 50000.0, response["personalDeduction"])
	}
}
