package repository

import (
	"github.com/CharlesSchiavinato/minsait-challenge-backend/model"
)

type CashLaunch interface {
	Insert(modelCashLaunch *model.CashLaunch) (*model.CashLaunch, error)
	List() (model.CashLaunches, error)
	GetByID(id int64) (*model.CashLaunch, error)
	Update(modelCashLaunch *model.CashLaunch) (*model.CashLaunch, error)
	DeleteByID(id int64) error
}
