package model

import "time"

type CashBalanceDaily struct {
	// Data de Referencia
	ReferenceDate time.Time `json:"reference_date" validate:"required" example:"2019-08-24T00:00:00Z" format:"date-time" minimum:"1900-01-01T00:00:00Z"`
	// Saldo
	Value float64 `json:"value" validate:"required" format:"float" example:"1.23" format:"float"`
}

type CashBalanceDailies []CashBalanceDaily

type CashBalanceDailyRangeReferenceDate struct {
	From time.Time
	To   time.Time
}
