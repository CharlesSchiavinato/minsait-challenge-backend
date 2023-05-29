package repository

import (
	"errors"
	"time"

	"github.com/CharlesSchiavinato/minsait-challenge-backend/model"
	"github.com/CharlesSchiavinato/minsait-challenge-backend/service/database/repository"
)

var cashLaunchIDLast int64 = 2

var InMemoryCashLaunches = model.CashLaunches{
	{
		ID:            1,
		ReferenceDate: time.Date(2001, 11, 22, 00, 00, 00, 000, time.UTC),
		Type:          "D",
		Description:   "Description InMemory 1",
		Value:         12.34,
		UpdatedAt:     time.Now().UTC(),
		CreatedAt:     time.Now().UTC(),
	},
	{
		ID:            2,
		ReferenceDate: time.Date(2000, 11, 22, 00, 00, 00, 000, time.UTC),
		Type:          "C",
		Description:   "Description InMemory 2",
		Value:         987.65,
		UpdatedAt:     time.Now().UTC(),
		CreatedAt:     time.Now().UTC(),
	},
	{
		ID:            3,
		ReferenceDate: time.Date(2000, 11, 22, 00, 00, 00, 000, time.UTC),
		Type:          "D",
		Description:   "Description InMemory 1",
		Value:         12.34,
		UpdatedAt:     time.Now().UTC(),
		CreatedAt:     time.Now().UTC(),
	},
}

type InMemoryCashLaunch struct {
	InMemory *InMemory
}

func NewCashLaunch(inMemory *InMemory) repository.CashLaunch {
	return &InMemoryCashLaunch{
		InMemory: inMemory,
	}
}

func (repositoryInMemoryCashLaunch *InMemoryCashLaunch) Insert(modelCashLaunch *model.CashLaunch) (*model.CashLaunch, error) {
	if repositoryInMemoryCashLaunch.InMemory.Error == true {
		return nil, errors.New("Error persist in database")
	}

	modelCashLaunchInsert := *modelCashLaunch
	cashLaunchIDLast += 1
	modelCashLaunchInsert.ID = cashLaunchIDLast
	InMemoryCashLaunches = append(InMemoryCashLaunches, modelCashLaunchInsert)

	return &modelCashLaunchInsert, nil
}

func (repositoryInMemoryCashLaunch *InMemoryCashLaunch) List() (model.CashLaunches, error) {
	if repositoryInMemoryCashLaunch.InMemory.Error == true {
		return nil, errors.New("Error load from database")
	}

	return InMemoryCashLaunches, nil
}

func (repositoryInMemoryCashLaunch *InMemoryCashLaunch) GetByID(id int64) (*model.CashLaunch, error) {
	if repositoryInMemoryCashLaunch.InMemory.Error == true {
		return nil, errors.New("Error load from database")
	}

	idx, modelCashLaunch := GetByID(id)

	if idx < 0 {
		return nil, repository.ErrNotFound{Message: "not found"}
	}

	return modelCashLaunch, nil
}

func (repositoryInMemoryCashLaunch *InMemoryCashLaunch) Update(modelCashLaunch *model.CashLaunch) (*model.CashLaunch, error) {
	if repositoryInMemoryCashLaunch.InMemory.Error == true {
		return nil, errors.New("Error persist in database")
	}

	idx, _ := GetByID(modelCashLaunch.ID)

	if idx < 0 {
		return nil, repository.ErrNotFound{Message: "not found"}
	}

	InMemoryCashLaunches[idx] = *modelCashLaunch

	return &InMemoryCashLaunches[idx], nil
}

func (repositoryInMemoryCashLaunch *InMemoryCashLaunch) DeleteByID(id int64) error {
	if repositoryInMemoryCashLaunch.InMemory.Error == true {
		return errors.New("Error persist in database")
	}

	idx, _ := GetByID(id)

	if idx < 0 {
		return repository.ErrNotFound{Message: "not found"}
	}

	return nil
}

func GetByID(id int64) (int, *model.CashLaunch) {
	for idx, cashLaunch := range InMemoryCashLaunches {
		if cashLaunch.ID == id {
			return idx, &cashLaunch
		}
	}

	return -1, nil
}
