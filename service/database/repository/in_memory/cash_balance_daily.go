package repository

import (
	"errors"
	"time"

	"github.com/CharlesSchiavinato/minsait-challenge-backend/model"
	"github.com/CharlesSchiavinato/minsait-challenge-backend/service/database/repository"
)

type InMemoryCashBalanceDaily struct {
	InMemory *InMemory
}

func NewCashBalanceDaily(inMemory *InMemory) repository.CashBalanceDaily {
	return &InMemoryCashBalanceDaily{
		InMemory: inMemory,
	}
}

func (repositoryInMemoryCashBalanceDaily *InMemoryCashBalanceDaily) GetByReferenceDate(referenceDate time.Time) (*model.CashBalanceDaily, error) {
	if repositoryInMemoryCashBalanceDaily.InMemory.Error == true {
		return nil, errors.New("Error load from database")
	}

	cashBalanceDaily := &model.CashBalanceDaily{
		ReferenceDate: referenceDate,
		Value:         0,
	}

	for _, cashLaunch := range InMemoryCashLaunches {
		if cashLaunch.ReferenceDate == referenceDate {
			if cashLaunch.Type == "C" {
				cashBalanceDaily.Value += cashLaunch.Value
			} else {
				cashBalanceDaily.Value -= cashLaunch.Value
			}
		}
	}

	return cashBalanceDaily, nil
}

func (repositoryInMemoryCashBalanceDaily *InMemoryCashBalanceDaily) GetByRangeReferenceDate(cashBalanceGetByRangeReferenceDateParams *model.CashBalanceDailyRangeReferenceDate) (model.CashBalanceDailies, error) {
	if repositoryInMemoryCashBalanceDaily.InMemory.Error == true {
		return nil, errors.New("Error load from database")
	}

	cashBalanceDailies := model.CashBalanceDailies{}

	for _, cashLaunch := range InMemoryCashLaunches {
		if cashLaunch.ReferenceDate.Sub(cashBalanceGetByRangeReferenceDateParams.From).Hours()/24 >= 0 &&
			cashBalanceGetByRangeReferenceDateParams.To.Sub(cashLaunch.ReferenceDate).Hours()/24 >= 0 {
			idx := getCashBalanceDailyByReferenceDate(cashBalanceDailies, cashLaunch.ReferenceDate)

			if idx < 0 {
				cashBalanceDailies = append(cashBalanceDailies, model.CashBalanceDaily{ReferenceDate: cashLaunch.ReferenceDate})
				idx = len(cashBalanceDailies) - 1
			}

			if cashLaunch.Type == "C" {
				cashBalanceDailies[idx].Value += cashLaunch.Value
			} else {
				cashBalanceDailies[idx].Value -= cashLaunch.Value
			}
		}
	}

	return cashBalanceDailies, nil
}

func getCashBalanceDailyByReferenceDate(cashBalanceDailies model.CashBalanceDailies, referenceDate time.Time) int {
	for index, cashBalanceDaily := range cashBalanceDailies {
		if cashBalanceDaily.ReferenceDate == referenceDate {
			return index
		}
	}

	return -1
}
