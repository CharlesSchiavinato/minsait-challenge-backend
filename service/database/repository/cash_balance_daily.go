package repository

import (
	"time"

	"github.com/CharlesSchiavinato/minsait-challenge-backend/model"
)

type CashBalanceDaily interface {
	GetByReferenceDate(referenceDate time.Time) (*model.CashBalanceDaily, error)
	GetByRangeReferenceDate(cashBalanceGetByRangeReferenceDateParams *model.CashBalanceDailyRangeReferenceDate) (model.CashBalanceDailies, error)
}
