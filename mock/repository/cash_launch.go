package mock_repository

import (
	"time"

	"github.com/CharlesSchiavinato/minsait-challenge-backend/model"
	"github.com/stretchr/testify/mock"
)

type MockCashLaunch struct {
	mock.Mock
}

func (mockCashLaunch *MockCashLaunch) Insert(modelCashLaunch *model.CashLaunch) (*model.CashLaunch, error) {
	args := mockCashLaunch.Called()

	modelCashLaunch = args.Get(0).(*model.CashLaunch)
	modelCashLaunchInsert := *modelCashLaunch

	modelCashLaunchInsert.ID = 1
	modelCashLaunchInsert.CreatedAt = time.Now().UTC()
	modelCashLaunchInsert.UpdatedAt = time.Now().UTC()

	return &modelCashLaunchInsert, args.Error(1)
}

func (mockCashLaunch *MockCashLaunch) List() (model.CashLaunches, error) {
	args := mockCashLaunch.Called()
	return args.Get(0).(model.CashLaunches), args.Error(1)
}

func (mockCashLaunch *MockCashLaunch) GetByID(id int64) (*model.CashLaunch, error) {
	args := mockCashLaunch.Called()
	return args.Get(0).(*model.CashLaunch), args.Error(1)
}

func (mockCashLaunch *MockCashLaunch) Update(modelCashLaunch *model.CashLaunch) (*model.CashLaunch, error) {
	return nil, nil
}

func (mockCashLaunch *MockCashLaunch) DeleteByID(id int64) error {
	return nil
}
