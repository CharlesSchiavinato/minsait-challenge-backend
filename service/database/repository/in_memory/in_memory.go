package repository

import (
	"github.com/CharlesSchiavinato/minsait-challenge-backend/service/database/repository"
)

type InMemory struct {
	Error bool
}

func NewInMemory(error bool) (repository.Repository, error) {
	return &InMemory{Error: error}, nil
}

func (inMemory *InMemory) Check() error {
	return nil
}

func (inMemory *InMemory) Close() error {
	return nil
}

func (inMemory *InMemory) CashLaunch() repository.CashLaunch {
	return NewCashLaunch(inMemory)
}

func (inMemory *InMemory) CashBalanceDaily() repository.CashBalanceDaily {
	return NewCashBalanceDaily(inMemory)
}
