package controller_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/CharlesSchiavinato/minsait-challenge-backend/controller"
	"github.com/CharlesSchiavinato/minsait-challenge-backend/model"
	repository "github.com/CharlesSchiavinato/minsait-challenge-backend/service/database/repository/postgres"
	"github.com/CharlesSchiavinato/minsait-challenge-backend/usecase"
	"github.com/CharlesSchiavinato/minsait-challenge-backend/util"
	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/assert"
)

var (
	config, _            = util.LoadConfig("./../")
	log                  = hclog.New(&hclog.LoggerOptions{Level: hclog.LevelFromString("OFF")})
	repositoryTest, _    = repository.NewPostgres(config)
	usecaseCashLaunch    = usecase.NewCashLaunch(repositoryTest.CashLaunch())
	controllerCashLaunch = controller.NewCashLaunch(log, usecaseCashLaunch)
	controllerTitle      = "CashLaunch"
)

func TestIntegrationCashLaunchInsert(t *testing.T) {
	type test struct {
		name                   string
		reqBodyModelCashLaunch *model.CashLaunch
		resBodyModel           interface{}
		repoError              bool
		wantResCode            int
		wantResBody            interface{}
		assert                 func(t *testing.T, reqBodyModelCashLaunch, resBodyModel *model.CashLaunch)
	}

	tests := []test{
		{
			name:         "BodyDeserializeError",
			resBodyModel: &model.Error{},
			wantResCode:  http.StatusBadRequest,
			wantResBody:  model.BadRequestDeserialize(controllerCashLaunchTitle),
		},
		{
			name:                   "ModelValidateError",
			reqBodyModelCashLaunch: &model.CashLaunch{},
			resBodyModel:           &model.Error{},
			wantResCode:            http.StatusBadRequest,
			wantResBody:            model.BadRequestModelValidate(controllerCashLaunchTitle, "The reference_date is empty;The type is empty;The description is empty;The value is less or equal 0"),
		},
		{
			name:                   "Success",
			reqBodyModelCashLaunch: modelCashLaunchDefault,
			resBodyModel:           &model.CashLaunch{},
			wantResCode:            http.StatusCreated,
			assert: func(t *testing.T, reqBodyModelCashLaunch, resBodyModel *model.CashLaunch) {
				usecase.CashLaunchModelFormat(reqBodyModelCashLaunch)

				assert.NotEqual(t, reqBodyModelCashLaunch.ID, resBodyModel.ID)
				assert.Equal(t, reqBodyModelCashLaunch.ReferenceDate, resBodyModel.ReferenceDate)
				assert.Equal(t, reqBodyModelCashLaunch.Type, resBodyModel.Type)
				assert.Equal(t, reqBodyModelCashLaunch.Description, resBodyModel.Description)
				assert.Equal(t, reqBodyModelCashLaunch.Value, resBodyModel.Value)
				assert.NotEqual(t, reqBodyModelCashLaunch.CreatedAt, resBodyModel.CreatedAt)
				assert.NotEqual(t, reqBodyModelCashLaunch.UpdatedAt, resBodyModel.UpdatedAt)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var bytesBody []byte

			if tt.reqBodyModelCashLaunch == nil {
				bytesBody = []byte("")
			} else {
				bytesBody, _ = json.Marshal(tt.reqBodyModelCashLaunch)
			}

			reqBody := bytes.NewBuffer(bytesBody)

			req, _ := http.NewRequest(http.MethodPost, "/api/cash/launch", reqBody)
			handler := http.HandlerFunc(controllerCashLaunch.Insert)
			res := httptest.NewRecorder()

			handler.ServeHTTP(res, req)

			if !reflect.DeepEqual(res.Code, tt.wantResCode) {
				t.Errorf("Insert() got res.code = %v, want %v", res.Code, tt.wantResCode)
			}

			json.NewDecoder(res.Body).Decode(&tt.resBodyModel)

			if tt.assert != nil {
				tt.assert(t, tt.reqBodyModelCashLaunch, tt.resBodyModel.(*model.CashLaunch))
			} else {
				if !reflect.DeepEqual(tt.resBodyModel, tt.wantResBody) {
					t.Errorf("Insert() got res.body = %v, want %v", tt.resBodyModel, tt.wantResBody)
				}
			}
		})
	}
}
