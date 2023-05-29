package usecase_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/CharlesSchiavinato/minsait-challenge-backend/model"
	repository_in_memory "github.com/CharlesSchiavinato/minsait-challenge-backend/service/database/repository/in_memory"
	"github.com/CharlesSchiavinato/minsait-challenge-backend/usecase"
)

func TestCashBalanceDailyGetByReferenceDate(t *testing.T) {
	type test struct {
		name                 string
		inputReferenceDate   time.Time
		wantCashBalanceDaily *model.CashBalanceDaily
		wantError            error
	}

	tests := []test{
		{
			name:               "ReferenceDateBeforeError",
			inputReferenceDate: usecase.CashLaunchReferenceDateMin.AddDate(0, 0, -1),
			wantError:          usecase.ErrParamValidate{Message: usecase.CashLaunchMessageReferenceDateBetweenError},
		},
		{
			name:               "ReferenceDateAfterError",
			inputReferenceDate: usecase.CashLaunchReferenceDateMax.AddDate(0, 0, 1),
			wantError:          usecase.ErrParamValidate{Message: usecase.CashLaunchMessageReferenceDateBetweenError},
		},
		{
			name:               "SuccessNotFound",
			inputReferenceDate: usecase.CashLaunchReferenceDateMin.AddDate(0, 0, 10),
			wantCashBalanceDaily: &model.CashBalanceDaily{
				ReferenceDate: usecase.CashLaunchReferenceDateMin.AddDate(0, 0, 10),
			},
			wantError: nil,
		},
		{
			name:               "SuccessFound",
			inputReferenceDate: time.Date(2000, 11, 22, 00, 00, 00, 000, time.UTC),
			wantCashBalanceDaily: &model.CashBalanceDaily{
				ReferenceDate: time.Date(2000, 11, 22, 00, 00, 00, 000, time.UTC),
				Value:         975.31,
			},
			wantError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository, _ := repository_in_memory.NewInMemory(false)
			repositoryCashBalanceDaily := repository.CashBalanceDaily()

			usecaseCashBalanceDaily := usecase.NewCashBalanceDaily(repositoryCashBalanceDaily)

			resultCashBalanceDaily, err := usecaseCashBalanceDaily.GetByReferenceDate(tt.inputReferenceDate)

			if !reflect.DeepEqual(err, tt.wantError) {
				t.Errorf("GetByReferenceDate() got error = %v, want = %v.", err, tt.wantError)
			}

			if !reflect.DeepEqual(resultCashBalanceDaily, tt.wantCashBalanceDaily) {
				t.Errorf("GetByReferenceDate() got result = %v, want = %v.", resultCashBalanceDaily, tt.wantCashBalanceDaily)
			}
		})
	}
}

func TestCashBalanceDailyGetByRangeReferenceDate(t *testing.T) {
	type test struct {
		name                                    string
		inputCashBalanceDailyRangeReferenceDate *model.CashBalanceDailyRangeReferenceDate
		wantCashBalanceDailies                  model.CashBalanceDailies
		wantError                               error
	}

	tests := []test{
		{
			name: "CashBalanceDailyRangeReferenceDateFromEmptyError",
			inputCashBalanceDailyRangeReferenceDate: &model.CashBalanceDailyRangeReferenceDate{
				To: usecase.CashLaunchReferenceDateMin,
			},
			wantError: usecase.ErrParamValidate{Message: usecase.CashBalanceDailyRangeReferenceDateFromEmptyError},
		},
		{
			name: "CashBalanceDailyRangeReferenceDateFromBeforeError",
			inputCashBalanceDailyRangeReferenceDate: &model.CashBalanceDailyRangeReferenceDate{
				From: usecase.CashLaunchReferenceDateMin.AddDate(0, 0, -1),
				To:   usecase.CashLaunchReferenceDateMin,
			},
			wantError: usecase.ErrParamValidate{Message: usecase.CashBalanceDailyRangeReferenceDateFromBetweenError},
		},
		{
			name: "CashBalanceDailyRangeReferenceDateFromAfterError",
			inputCashBalanceDailyRangeReferenceDate: &model.CashBalanceDailyRangeReferenceDate{
				From: usecase.CashLaunchReferenceDateMax.AddDate(0, 0, 1),
				To:   usecase.CashLaunchReferenceDateMin,
			},
			wantError: usecase.ErrParamValidate{Message: usecase.CashBalanceDailyRangeReferenceDateFromBetweenError},
		},
		{
			name: "CashBalanceDailyRangeReferenceDateToEmptyError",
			inputCashBalanceDailyRangeReferenceDate: &model.CashBalanceDailyRangeReferenceDate{
				From: usecase.CashLaunchReferenceDateMin,
			},
			wantError: usecase.ErrParamValidate{Message: usecase.CashBalanceDailyRangeReferenceDateToEmptyError},
		},
		{
			name: "CashBalanceDailyRangeReferenceDateToBeforeError",
			inputCashBalanceDailyRangeReferenceDate: &model.CashBalanceDailyRangeReferenceDate{
				From: usecase.CashLaunchReferenceDateMin,
				To:   usecase.CashLaunchReferenceDateMin.AddDate(0, 0, -1),
			},
			wantError: usecase.ErrParamValidate{Message: usecase.CashBalanceDailyRangeReferenceDateToBetweenError},
		},
		{
			name: "CashBalanceDailyRangeReferenceDateToAfterError",
			inputCashBalanceDailyRangeReferenceDate: &model.CashBalanceDailyRangeReferenceDate{
				From: usecase.CashLaunchReferenceDateMax,
				To:   usecase.CashLaunchReferenceDateMax.AddDate(0, 0, 1),
			},
			wantError: usecase.ErrParamValidate{Message: usecase.CashBalanceDailyRangeReferenceDateToBetweenError},
		},
		{
			name: "CashBalanceDailyRangeReferenceDateToSmallerFromError",
			inputCashBalanceDailyRangeReferenceDate: &model.CashBalanceDailyRangeReferenceDate{
				From: usecase.CashLaunchReferenceDateMin.AddDate(0, 0, 2),
				To:   usecase.CashLaunchReferenceDateMin.AddDate(0, 0, 1),
			},
			wantError: usecase.ErrParamValidate{Message: usecase.CashBalanceDailyRangeReferenceDateToSmallerFromError},
		},
		{
			name: "CashBalanceDailyRangeReferenceDateRangeError",
			inputCashBalanceDailyRangeReferenceDate: &model.CashBalanceDailyRangeReferenceDate{
				From: usecase.CashLaunchReferenceDateMin,
				To:   usecase.CashLaunchReferenceDateMin.AddDate(0, 0, 32),
			},
			wantError: usecase.ErrParamValidate{Message: usecase.CashBalanceDailyRangeReferenceDateRangeError},
		},
		{
			name: "Success",
			inputCashBalanceDailyRangeReferenceDate: &model.CashBalanceDailyRangeReferenceDate{
				From: time.Date(2000, 11, 22, 00, 00, 00, 000, time.UTC),
				To:   time.Date(2000, 11, 22, 00, 00, 00, 000, time.UTC),
			},
			wantCashBalanceDailies: model.CashBalanceDailies{
				{
					ReferenceDate: time.Date(2000, 11, 22, 00, 00, 00, 000, time.UTC),
					Value:         975.31,
				},
			},
			wantError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository, _ := repository_in_memory.NewInMemory(false)
			repositoryCashBalanceDaily := repository.CashBalanceDaily()

			usecaseCashBalanceDaily := usecase.NewCashBalanceDaily(repositoryCashBalanceDaily)

			resultCashBalanceDailies, err := usecaseCashBalanceDaily.GetByRangeReferenceDate(tt.inputCashBalanceDailyRangeReferenceDate)

			if !reflect.DeepEqual(err, tt.wantError) {
				t.Errorf("GetByReferenceDate() got error = %v, want = %v.", err, tt.wantError)
			}

			if !reflect.DeepEqual(resultCashBalanceDailies, tt.wantCashBalanceDailies) {
				t.Errorf("GetByReferenceDate() got result = %v, want = %v.", resultCashBalanceDailies, tt.wantCashBalanceDailies)
			}
		})
	}
}
