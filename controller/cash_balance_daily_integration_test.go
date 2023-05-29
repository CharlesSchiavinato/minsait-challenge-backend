package controller_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/CharlesSchiavinato/minsait-challenge-backend/controller"
	"github.com/CharlesSchiavinato/minsait-challenge-backend/model"
	"github.com/CharlesSchiavinato/minsait-challenge-backend/usecase"
)

var (
	usecaseCashBalanceDaily         = usecase.NewCashBalanceDaily(repositoryTest.CashBalanceDaily())
	controllerCashBalanceDaily      = controller.NewCashBalanceDaily(log, usecaseCashBalanceDaily)
	controllerCashBalanceDailyTitle = "CashBalanceDaily"
)

func TestIntegrationCashLaunchBalanceDailyByReferenceDate(t *testing.T) {
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
			name:         "Success",
			reqParam:     usecase.CashLaunchReferenceDateMin.AddDate(0, 0, 11).Format("2006-01-02"),
			resBodyModel: &model.CashBalanceDaily{},
			wantResCode:  http.StatusOK,
			wantResBody:  &model.CashBalanceDaily{ReferenceDate: usecase.CashLaunchReferenceDateMin.AddDate(0, 0, 11)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
