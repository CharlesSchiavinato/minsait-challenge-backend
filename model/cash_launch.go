package model

import "time"

type CashLaunch struct {
	// Identificador do Lançamento (Gerado automaticamente na inclusão)
	ID int64 `json:"id" validate:"required" minimum:"1" format:"int64"`
	// Data de Referencia do Lançamento
	ReferenceDate time.Time `json:"reference_date" validate:"required" example:"2019-08-24T00:00:00Z" format:"date-time" minimum:"1900-01-01T00:00:00Z"`
	// Tipo do Lançamento (C=Crédito D=Débito)
	Type string `json:"type" validate:"required" enums:"C,D"`
	// Descrição do Lançamento
	Description string `json:"description" validate:"required"`
	// Valor do Lançamento
	Value float64 `json:"value" validate:"required" example:"1.23" format:"float"`
	// Data da Última Alteração do Lançamento (Atualizado automaticamente na inclusão e alteração)
	UpdatedAt time.Time `json:"updated_at" validate:"required" example:"2019-08-24T16:59:59Z" format:"date-time"`
	// Data de Inclusão do Lançamento (Gerado automaticamente na inclusão)
	CreatedAt time.Time `json:"created_at" validate:"required" example:"2019-08-24T16:59:59Z" format:"date-time"`
}

type CashLaunches []CashLaunch

type parametersCashLaunchWrapper struct {
	// Data de Referencia do Lançamento
	ReferenceDate time.Time `json:"reference_date" validate:"required" example:"2019-08-24T00:00:00Z" format:"date-time" minimum:"1900-01-01T00:00:00Z"`
	// Tipo do Lançamento (C=Crédito D=Débito)
	Type string `json:"type" validate:"required" enums:"C,D"`
	// Descrição do Lançamento
	Description string `json:"description" validate:"required"`
	// Valor do Lançamento
	Value float64 `json:"value" validate:"required" example:"1.23" format:"float"`
}
