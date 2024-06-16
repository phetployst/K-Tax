package models

import "gorm.io/gorm"

type Allowance struct {
	gorm.Model
	AllowanceType    string  `json:"allowanceType"`
	Amount           float64 `json:"amount"`
	TaxCalculationID uint    `json:"-"`
}

type TaxCalculation struct {
	gorm.Model
	TotalIncome float64     `json:"totalIncome"`
	WHT         float64     `json:"wht"`
	Allowances  []Allowance `json:"allowances" gorm:"foreignkey:TaxCalculationID"`
	Tax         float64     `json:"tax"`
	TaxRefund   float64     `json:"taxRefund"`
}

type AdminSetting struct {
	gorm.Model
	PersonalDeduction float64 `json:"personalDeduction"`
	KReceiptDeduction float64 `json:"kReceiptDeduction"`
}
