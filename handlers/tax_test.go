package handlers

import (
	"testing"
)

func TestCalculateTaxAmount(t *testing.T) {
}

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
