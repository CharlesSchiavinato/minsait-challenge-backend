package controller_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/CharlesSchiavinato/minsait-challenge-backend/controller"
	"github.com/CharlesSchiavinato/minsait-challenge-backend/model"
	repository_in_memory "github.com/CharlesSchiavinato/minsait-challenge-backend/service/database/repository/in_memory"
	"github.com/CharlesSchiavinato/minsait-challenge-backend/usecase"
	"github.com/hashicorp/go-hclog"
)

func TestCashBalanceDailyGetByReferenceDate(t *testing.T) {
	type test struct {
		name         string
		reqParam     string
		resBodyModel interface{}
		repoError    bool
		wantResCode  int
		wantResBody  interface{}
	}

	tests := []test{
		{
			name:         "ParamInvalidError",
			reqParam:     "x",
			resBodyModel: &model.Error{},
			wantResCode:  http.StatusBadRequest,
			wantResBody:  model.BadRequestParamValidate("Date invalid"),
		},
		{
			name:         "ParamBetweenError",
			reqParam:     usecase.CashLaunchReferenceDateMin.AddDate(0, 0, -1).Format("2006-01-02"),
			resBodyModel: &model.Error{},
			wantResCode:  http.StatusBadRequest,
			wantResBody:  model.BadRequestParamValidate(usecase.CashLaunchMessageReferenceDateBetweenError),
		},
		{
			name:         "RepositoryError",
			reqParam:     usecase.CashLaunchReferenceDateMin.Format("2006-01-02"),
			repoError:    true,
			resBodyModel: &model.Error{},
			wantResCode:  http.StatusInternalServerError,
			wantResBody:  model.InternalServerErrorRepositoryLoad(controllerCashBalanceDailyTitle),
		},
		{
			name:         "Success",
			reqParam:     usecase.CashLaunchReferenceDateMin.AddDate(0, 0, 11).Format("2006-01-02"),
			resBodyModel: &model.CashBalanceDaily{},
			wantResCode:  http.StatusOK,
			wantResBody:  &model.CashBalanceDaily{ReferenceDate: usecase.CashLaunchReferenceDateMin.AddDate(0, 0, 11)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log := hclog.New(&hclog.LoggerOptions{Level: hclog.LevelFromString("OFF")})
			repository, _ := repository_in_memory.NewInMemory(tt.repoError)
			repositoryCashBalanceDaily := repository.CashBalanceDaily()
			usecaseCashBalanceDaily := usecase.NewCashBalanceDaily(repositoryCashBalanceDaily)

			controllerCashBalanceDaily := controller.NewCashBalanceDaily(log, usecaseCashBalanceDaily)

			url := fmt.Sprintf("/api/cash/balance/daily/%v", tt.reqParam)

			req, _ := http.NewRequest(http.MethodGet, url, nil)
			handler := http.HandlerFunc(controllerCashBalanceDaily.GetByReferenceDate)
			res := httptest.NewRecorder()

			handler.ServeHTTP(res, req)

			if !reflect.DeepEqual(res.Code, tt.wantResCode) {
				t.Errorf("GetByReferenceDate() got res.code = %v, want %v", res.Code, tt.wantResCode)
			}

			json.NewDecoder(res.Body).Decode(tt.resBodyModel)

			if !reflect.DeepEqual(tt.resBodyModel, tt.wantResBody) {
				t.Errorf("GetByReferenceDate() got res.body = %v, want %v", tt.resBodyModel, tt.wantResBody)
			}
		})
	}
}

func TestCashBalanceDailyGetByRangeReferenceDate(t *testing.T) {
	type test struct {
		name         string
		reqParam     string
		resBodyModel interface{}
		repoError    bool
		wantResCode  int
		wantResBody  interface{}
	}

	tests := []test{
		{
			name:         "ParamEmptyError",
			reqParam:     "",
			resBodyModel: &model.Error{},
			wantResCode:  http.StatusBadRequest,
			wantResBody:  model.BadRequestParamValidate("The param from is empty;The param to is empty"),
		},
		{
			name:         "ParamInvalidError",
			reqParam:     "?from=2020-01-02&to=2020-01-01",
			resBodyModel: &model.Error{},
			wantResCode:  http.StatusBadRequest,
			wantResBody:  model.BadRequestParamValidate("The param to is smaller the param from"),
		},
		{
			name:         "RepositoryError",
			reqParam:     "?from=2020-01-01&to=2020-01-01",
			repoError:    true,
			resBodyModel: &model.Error{},
			wantResCode:  http.StatusInternalServerError,
			wantResBody:  model.InternalServerErrorRepositoryLoad(controllerCashBalanceDailyTitle),
		},
		{
			name:         "Success",
			reqParam:     "?from=2000-11-22&to=2000-11-22",
			resBodyModel: &model.CashBalanceDailies{},
			wantResCode:  http.StatusOK,
			wantResBody: &model.CashBalanceDailies{
				{
					ReferenceDate: time.Date(2000, 11, 22, 00, 00, 00, 000, time.UTC),
					Value:         975.31,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log := hclog.New(&hclog.LoggerOptions{Level: hclog.LevelFromString("OFF")})
			repository, _ := repository_in_memory.NewInMemory(tt.repoError)
			repositoryCashBalanceDaily := repository.CashBalanceDaily()
			usecaseCashBalanceDaily := usecase.NewCashBalanceDaily(repositoryCashBalanceDaily)

			controllerCashBalanceDaily := controller.NewCashBalanceDaily(log, usecaseCashBalanceDaily)

			url := fmt.Sprintf("/api/cash/balance/daily%v", tt.reqParam)

			req, _ := http.NewRequest(http.MethodGet, url, nil)
			handler := http.HandlerFunc(controllerCashBalanceDaily.GetByRangeReferenceDate)
			res := httptest.NewRecorder()

			handler.ServeHTTP(res, req)

			if !reflect.DeepEqual(res.Code, tt.wantResCode) {
				t.Errorf("GetByReferenceDate() got res.code = %v, want %v", res.Code, tt.wantResCode)
			}

			json.NewDecoder(res.Body).Decode(tt.resBodyModel)

			if !reflect.DeepEqual(tt.resBodyModel, tt.wantResBody) {
				t.Errorf("GetByReferenceDate() got res.body = %v, want %v", tt.resBodyModel, tt.wantResBody)
			}
		})
	}
}
