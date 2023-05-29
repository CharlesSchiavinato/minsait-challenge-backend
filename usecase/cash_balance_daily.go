package usecase

import (
	"fmt"
	"strings"
	"time"

	"github.com/CharlesSchiavinato/minsait-challenge-backend/model"
	"github.com/CharlesSchiavinato/minsait-challenge-backend/service/database/repository"
	"github.com/CharlesSchiavinato/minsait-challenge-backend/util"
)

var (
	CashBalanceDailyRangeReferenceDateFromEmptyError     = "The param from is empty"
	CashBalanceDailyRangeReferenceDateFromBetweenError   = fmt.Sprintf("The param from value is not between %v and %v", CashLaunchReferenceDateMin, CashLaunchReferenceDateMax)
	CashBalanceDailyRangeReferenceDateToEmptyError       = "The param to is empty"
	CashBalanceDailyRangeReferenceDateToBetweenError     = fmt.Sprintf("The param to value is not between %v and %v", CashLaunchReferenceDateMin, CashLaunchReferenceDateMax)
	CashBalanceDailyRangeReferenceDateToSmallerFromError = "The param to is smaller the param from"
	CashBalanceDailyRangeReferenceDateRangeError         = "the range is greater than 31 days"
)

type CashBalanceDaily interface {
	GetByReferenceDate(referenceDate time.Time) (*model.CashBalanceDaily, error)
	GetByRangeReferenceDate(cashBalanceGetByRangeReferenceDateParams *model.CashBalanceDailyRangeReferenceDate) (model.CashBalanceDailies, error)
}

type UseCaseCashBalanceDaily struct {
	RepositoryCashBalanceDaily repository.CashBalanceDaily
}

func NewCashBalanceDaily(repositoryCashBalanceDaily repository.CashBalanceDaily) CashBalanceDaily {
	return &UseCaseCashBalanceDaily{
		RepositoryCashBalanceDaily: repositoryCashBalanceDaily,
	}
}

func (useCaseCashBalanceDaily *UseCaseCashBalanceDaily) GetByReferenceDate(referenceDate time.Time) (*model.CashBalanceDaily, error) {
	err := cashLaunchReferenceDateValidate(referenceDate)

	if err != nil {
		return nil, err
	}

	modelCashBalanceDaily, err := useCaseCashBalanceDaily.RepositoryCashBalanceDaily.GetByReferenceDate(referenceDate)

	if err != nil {
		if _, ok := err.(repository.ErrNotFound); ok {
			return &model.CashBalanceDaily{ReferenceDate: referenceDate}, nil
		} else {
			return nil, err
		}
	}

	modelCashBalanceDaily.Value = util.MathRoundPrecision(modelCashBalanceDaily.Value, 2)

	return modelCashBalanceDaily, err
}

func (useCaseCashBalanceDaily *UseCaseCashBalanceDaily) GetByRangeReferenceDate(cashBalanceGetByRangeReferenceDateParams *model.CashBalanceDailyRangeReferenceDate) (model.CashBalanceDailies, error) {
	err := CashBalanceDailyRangeReferenceDateValidate(cashBalanceGetByRangeReferenceDateParams)

	if err != nil {
		return nil, err
	}

	return useCaseCashBalanceDaily.RepositoryCashBalanceDaily.GetByRangeReferenceDate(cashBalanceGetByRangeReferenceDateParams)
}

func CashBalanceDailyRangeReferenceDateValidate(cashBalanceGetByRangeReferenceDateParams *model.CashBalanceDailyRangeReferenceDate) error {
	messages := []string{}

	if cashBalanceGetByRangeReferenceDateParams.From.IsZero() {
		messages = append(messages, CashBalanceDailyRangeReferenceDateFromEmptyError)
	} else if cashBalanceGetByRangeReferenceDateParams.From.Before(CashLaunchReferenceDateMin) ||
		cashBalanceGetByRangeReferenceDateParams.From.After(CashLaunchReferenceDateMax) {
		messages = append(messages, CashBalanceDailyRangeReferenceDateFromBetweenError)
	}

	if cashBalanceGetByRangeReferenceDateParams.To.IsZero() {
		messages = append(messages, CashBalanceDailyRangeReferenceDateToEmptyError)
	} else if cashBalanceGetByRangeReferenceDateParams.To.Before(CashLaunchReferenceDateMin) ||
		cashBalanceGetByRangeReferenceDateParams.To.After(CashLaunchReferenceDateMax) {
		messages = append(messages, CashBalanceDailyRangeReferenceDateToBetweenError)
	}

	if len(messages) > 0 {
		return ErrParamValidate{Message: strings.Join(messages, ";")}
	}

	dateDiffDays := int64(cashBalanceGetByRangeReferenceDateParams.To.Sub(cashBalanceGetByRangeReferenceDateParams.From).Hours() / 24)

	if dateDiffDays < 0 {
		messages = append(messages, CashBalanceDailyRangeReferenceDateToSmallerFromError)
	} else if dateDiffDays > 31 {
		messages = append(messages, CashBalanceDailyRangeReferenceDateRangeError)
	}

	if len(messages) > 0 {
		return ErrParamValidate{Message: strings.Join(messages, ";")}
	}

	return nil
}
