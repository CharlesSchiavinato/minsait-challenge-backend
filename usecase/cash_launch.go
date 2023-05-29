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
	CashLaunchReferenceDateMin  = time.Date(1900, 01, 01, 00, 00, 00, 000, time.UTC)
	CashLaunchReferenceDateMax  = time.Now().UTC().AddDate(10, 0, 0)
	CashLaunchDescriptionMinLen = 3
	CashLaunchDescriptionMaxLen = 100

	CashLaunchMessageReferenceDateEmptyError   = "The reference_date is empty"
	CashLaunchMessageReferenceDateBetweenError = fmt.Sprintf("The reference_date value is not between %v and %v", CashLaunchReferenceDateMin, CashLaunchReferenceDateMax)
	CashLaunchMessageTypeEmptyError            = "The type is empty"
	CashLaunchMessageTypeInvalidError          = "The type not in ['C', 'D']"
	CashLaunchMessageDescriptionEmptyError     = "The description is empty"
	CashLaunchMessageDescriptionSizeError      = fmt.Sprintf("The description size is not between %v and %v", CashLaunchDescriptionMinLen, CashLaunchDescriptionMaxLen)
	CashLaunchMessageValueError                = "The value is less or equal 0"
)

type CashLaunch interface {
	Insert(modelCashLaunch *model.CashLaunch) (*model.CashLaunch, error)
	List() (model.CashLaunches, error)
	GetByID(id int64) (*model.CashLaunch, error)
	Update(modelCashLaunch *model.CashLaunch) (*model.CashLaunch, error)
	DeleteByID(id int64) error
}

type UseCaseCashLaunch struct {
	RepositoryCashLaunch repository.CashLaunch
}

func NewCashLaunch(repositoryCashLaunch repository.CashLaunch) CashLaunch {
	return &UseCaseCashLaunch{
		RepositoryCashLaunch: repositoryCashLaunch,
	}
}

func (useCaseCashLaunch *UseCaseCashLaunch) Insert(modelCashLaunch *model.CashLaunch) (*model.CashLaunch, error) {
	err := cashLaunchModelValidate(modelCashLaunch)

	if err != nil {
		return nil, err
	}

	modelCashLaunch.CreatedAt = time.Now().UTC()
	modelCashLaunch.UpdatedAt = modelCashLaunch.CreatedAt

	return useCaseCashLaunch.RepositoryCashLaunch.Insert(modelCashLaunch)
}

func (useCaseCashLaunch *UseCaseCashLaunch) List() (model.CashLaunches, error) {
	return useCaseCashLaunch.RepositoryCashLaunch.List()
}

func (useCaseCashLaunch *UseCaseCashLaunch) GetByID(id int64) (*model.CashLaunch, error) {
	return useCaseCashLaunch.RepositoryCashLaunch.GetByID(id)
}

func (useCaseCashLaunch *UseCaseCashLaunch) Update(modelCashLaunch *model.CashLaunch) (*model.CashLaunch, error) {
	err := cashLaunchModelValidate(modelCashLaunch)

	if err != nil {
		return nil, err
	}

	modelCashLaunch.UpdatedAt = time.Now().UTC()

	return useCaseCashLaunch.RepositoryCashLaunch.Update(modelCashLaunch)
}

func (useCaseCashLaunch *UseCaseCashLaunch) DeleteByID(id int64) error {
	return useCaseCashLaunch.RepositoryCashLaunch.DeleteByID(id)
}

func cashLaunchModelValidate(modelCashLaunch *model.CashLaunch) error {
	messages := []string{}

	CashLaunchModelFormat(modelCashLaunch)

	err := cashLaunchReferenceDateValidate(modelCashLaunch.ReferenceDate)

	if err != nil {
		messages = append(messages, err.Error())
	}

	if modelCashLaunch.Type == "" {
		messages = append(messages, CashLaunchMessageTypeEmptyError)
	} else if modelCashLaunch.Type != "C" && modelCashLaunch.Type != "D" {
		messages = append(messages, CashLaunchMessageTypeInvalidError)
	}

	if modelCashLaunch.Description == "" {
		messages = append(messages, CashLaunchMessageDescriptionEmptyError)
	} else if len(modelCashLaunch.Description) < CashLaunchDescriptionMinLen ||
		len(modelCashLaunch.Description) > CashLaunchDescriptionMaxLen {
		messages = append(messages, CashLaunchMessageDescriptionSizeError)
	}

	if modelCashLaunch.Value <= 0 {
		messages = append(messages, CashLaunchMessageValueError)
	}

	if len(messages) > 0 {
		return ErrModelValidate{Message: strings.Join(messages, ";")}
	}

	return nil
}

func cashLaunchReferenceDateValidate(referenceDate time.Time) error {
	messages := []string{}

	if referenceDate.IsZero() {
		messages = append(messages, CashLaunchMessageReferenceDateEmptyError)
	} else if referenceDate.Before(CashLaunchReferenceDateMin) ||
		referenceDate.After(CashLaunchReferenceDateMax) {
		messages = append(messages, CashLaunchMessageReferenceDateBetweenError)
	}

	if len(messages) > 0 {
		return ErrParamValidate{Message: strings.Join(messages, ";")}
	}

	return nil
}

func CashLaunchModelFormat(modelCashLaunch *model.CashLaunch) {
	modelCashLaunch.Description = util.FormatTitle(modelCashLaunch.Description)
	modelCashLaunch.Type = util.FormatTextWithoutSpace(util.FormatTitle(modelCashLaunch.Type))
	modelCashLaunch.Value = util.MathRoundPrecision(modelCashLaunch.Value, 2)
}
