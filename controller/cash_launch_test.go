package controller_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/CharlesSchiavinato/minsait-challenge-backend/controller"
	"github.com/CharlesSchiavinato/minsait-challenge-backend/model"
	repository_in_memory "github.com/CharlesSchiavinato/minsait-challenge-backend/service/database/repository/in_memory"
	"github.com/CharlesSchiavinato/minsait-challenge-backend/usecase"
	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/assert"
)

var controllerCashLaunchTitle = "CashLaunch"
var modelCashLaunchDefault = &model.CashLaunch{
	ReferenceDate: usecase.CashLaunchReferenceDateMin,
	Type:          "c",
	Description:   "Description Test",
	Value:         0.01,
}

func TestCashLaunchInsert(t *testing.T) {
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
			name:                   "RepositoryError",
			reqBodyModelCashLaunch: modelCashLaunchDefault,
			repoError:              true,
			resBodyModel:           &model.Error{},
			wantResCode:            http.StatusInternalServerError,
			wantResBody:            model.InternalServerErrorRepositoryPersist(controllerCashLaunchTitle),
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
			log := hclog.New(&hclog.LoggerOptions{Level: hclog.LevelFromString("OFF")})
			repository, _ := repository_in_memory.NewInMemory(tt.repoError)
			repositoryCashLaunch := repository.CashLaunch()
			usecaseCashLaunch := usecase.NewCashLaunch(repositoryCashLaunch)
			controllerCashLaunch := controller.NewCashLaunch(log, usecaseCashLaunch)

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

func TestCashLaunchList(t *testing.T) {
	type test struct {
		name         string
		resBodyModel interface{}
		repoError    bool
		wantResCode  int
		wantResBody  interface{}
	}

	tests := []test{
		{
			name:         "RepositoryError",
			repoError:    true,
			resBodyModel: &model.Error{},
			wantResCode:  http.StatusInternalServerError,
			wantResBody:  model.InternalServerErrorRepositoryLoad(controllerCashLaunchTitle),
		},
		{
			name:         "Success",
			resBodyModel: &model.CashLaunches{},
			wantResCode:  http.StatusOK,
			wantResBody:  &repository_in_memory.InMemoryCashLaunches,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log := hclog.New(&hclog.LoggerOptions{Level: hclog.LevelFromString("OFF")})
			repository, _ := repository_in_memory.NewInMemory(tt.repoError)
			repositoryCashLaunch := repository.CashLaunch()
			usecaseCashLaunch := usecase.NewCashLaunch(repositoryCashLaunch)
			controllerCashLaunch := controller.NewCashLaunch(log, usecaseCashLaunch)

			req, _ := http.NewRequest(http.MethodGet, "/api/cash/launch", nil)
			handler := http.HandlerFunc(controllerCashLaunch.List)
			res := httptest.NewRecorder()

			handler.ServeHTTP(res, req)

			if !reflect.DeepEqual(res.Code, tt.wantResCode) {
				t.Errorf("List() got res.code = %v, want %v", res.Code, tt.wantResCode)
			}

			json.NewDecoder(res.Body).Decode(&tt.resBodyModel)

			if !reflect.DeepEqual(tt.resBodyModel, tt.wantResBody) {
				t.Errorf("List() got res.body = %v, want %v", tt.resBodyModel, tt.wantResBody)
			}
		})
	}
}

func TestCashLaunchGetByID(t *testing.T) {
	_, modelCashLaunchRes := repository_in_memory.GetByID(1)

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
			name:         "ParamError",
			reqParam:     "x",
			resBodyModel: &model.Error{},
			wantResCode:  http.StatusBadRequest,
			wantResBody:  model.BadRequestParamValidate("Id invalid"),
		},
		{
			name:         "NotFoundError",
			reqParam:     "0",
			resBodyModel: &model.Error{},
			wantResCode:  http.StatusNotFound,
			wantResBody:  model.NotFound(controllerCashLaunchTitle),
		},
		{
			name:         "RepositoryError",
			reqParam:     "1",
			repoError:    true,
			resBodyModel: &model.Error{},
			wantResCode:  http.StatusInternalServerError,
			wantResBody:  model.InternalServerErrorRepositoryLoad(controllerCashLaunchTitle),
		},
		{
			name:         "Success",
			reqParam:     "1",
			resBodyModel: &model.CashLaunch{},
			wantResCode:  http.StatusOK,
			wantResBody:  modelCashLaunchRes,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log := hclog.New(&hclog.LoggerOptions{Level: hclog.LevelFromString("OFF")})
			repository, _ := repository_in_memory.NewInMemory(tt.repoError)
			repositoryCashLaunch := repository.CashLaunch()
			usecaseCashLaunch := usecase.NewCashLaunch(repositoryCashLaunch)
			controllerCashLaunch := controller.NewCashLaunch(log, usecaseCashLaunch)

			url := fmt.Sprintf("/api/cash/launch/%v", tt.reqParam)

			req, _ := http.NewRequest(http.MethodGet, url, nil)
			handler := http.HandlerFunc(controllerCashLaunch.GetByID)
			res := httptest.NewRecorder()

			handler.ServeHTTP(res, req)

			if !reflect.DeepEqual(res.Code, tt.wantResCode) {
				t.Errorf("GetByID() got res.code = %v, want %v", res.Code, tt.wantResCode)
			}

			json.NewDecoder(res.Body).Decode(&tt.resBodyModel)

			if !reflect.DeepEqual(tt.resBodyModel, tt.wantResBody) {
				t.Errorf("GetByID() got res.body = %v, want %v", tt.resBodyModel, tt.wantResBody)
			}
		})
	}
}

func TestCashLaunchUpdate(t *testing.T) {
	_, modelCashLaunchRes := repository_in_memory.GetByID(1)

	modelCashLaunchRes.Description = fmt.Sprintf("%v Test Update", modelCashLaunchRes.Description)
	modelCashLaunchRes.ReferenceDate = modelCashLaunchDefault.ReferenceDate.AddDate(1, 1, 1)
	modelCashLaunchRes.Value += 98.76

	type test struct {
		name                   string
		reqParam               string
		reqBodyModelCashLaunch *model.CashLaunch
		resBodyModel           interface{}
		repoError              bool
		wantResCode            int
		wantResBody            interface{}
		assert                 func(t *testing.T, reqBodyModelCashLaunch, resBodyModel *model.CashLaunch)
	}

	tests := []test{
		{
			name:         "ParamError",
			reqParam:     "x",
			resBodyModel: &model.Error{},
			wantResCode:  http.StatusBadRequest,
			wantResBody:  model.BadRequestParamValidate("Id invalid"),
		},
		{
			name:         "BodyDeserializeError",
			reqParam:     "1",
			resBodyModel: &model.Error{},
			wantResCode:  http.StatusBadRequest,
			wantResBody:  model.BadRequestDeserialize(controllerCashLaunchTitle),
		},
		{
			name:                   "ModelValidateError",
			reqParam:               "1",
			reqBodyModelCashLaunch: &model.CashLaunch{},
			resBodyModel:           &model.Error{},
			wantResCode:            http.StatusBadRequest,
			wantResBody:            model.BadRequestModelValidate(controllerCashLaunchTitle, "The reference_date is empty;The type is empty;The description is empty;The value is less or equal 0"),
		},
		{
			name:                   "NotFoundError",
			reqParam:               "0",
			reqBodyModelCashLaunch: modelCashLaunchDefault,
			resBodyModel:           &model.Error{},
			wantResCode:            http.StatusNotFound,
			wantResBody:            model.NotFound(controllerCashLaunchTitle),
		},
		{
			name:                   "RepositoryError",
			reqParam:               "1",
			reqBodyModelCashLaunch: modelCashLaunchDefault,
			repoError:              true,
			resBodyModel:           &model.Error{},
			wantResCode:            http.StatusInternalServerError,
			wantResBody:            model.InternalServerErrorRepositoryPersist(controllerCashLaunchTitle),
		},
		{
			name:                   "Success",
			reqParam:               "1",
			reqBodyModelCashLaunch: modelCashLaunchRes,
			resBodyModel:           &model.CashLaunch{},
			wantResCode:            http.StatusOK,
			assert: func(t *testing.T, reqBodyModelCashLaunch, resBodyModel *model.CashLaunch) {
				usecase.CashLaunchModelFormat(reqBodyModelCashLaunch)

				assert.Equal(t, reqBodyModelCashLaunch.ID, resBodyModel.ID)
				assert.Equal(t, reqBodyModelCashLaunch.ReferenceDate, resBodyModel.ReferenceDate)
				assert.Equal(t, reqBodyModelCashLaunch.Description, resBodyModel.Description)
				assert.Equal(t, reqBodyModelCashLaunch.Value, resBodyModel.Value)
				assert.Equal(t, reqBodyModelCashLaunch.CreatedAt, resBodyModel.CreatedAt)
				assert.NotEqual(t, reqBodyModelCashLaunch.UpdatedAt, resBodyModel.UpdatedAt)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log := hclog.New(&hclog.LoggerOptions{Level: hclog.LevelFromString("OFF")})
			repository, _ := repository_in_memory.NewInMemory(tt.repoError)
			repositoryCashLaunch := repository.CashLaunch()
			usecaseCashLaunch := usecase.NewCashLaunch(repositoryCashLaunch)
			controllerCashLaunch := controller.NewCashLaunch(log, usecaseCashLaunch)

			url := fmt.Sprintf("/api/cash/launch/%v", tt.reqParam)

			var bytesBody []byte

			if tt.reqBodyModelCashLaunch == nil {
				bytesBody = []byte("")
			} else {
				bytesBody, _ = json.Marshal(tt.reqBodyModelCashLaunch)
			}

			reqBody := bytes.NewBuffer(bytesBody)

			req, _ := http.NewRequest(http.MethodPut, url, reqBody)
			handler := http.HandlerFunc(controllerCashLaunch.Update)
			res := httptest.NewRecorder()

			handler.ServeHTTP(res, req)

			if !reflect.DeepEqual(res.Code, tt.wantResCode) {
				t.Errorf("Update() got res.code = %v, want %v", res.Code, tt.wantResCode)
			}

			json.NewDecoder(res.Body).Decode(&tt.resBodyModel)

			if tt.assert != nil {
				tt.assert(t, tt.reqBodyModelCashLaunch, tt.resBodyModel.(*model.CashLaunch))
			} else {
				if !reflect.DeepEqual(tt.resBodyModel, tt.wantResBody) {
					t.Errorf("Update() got res.body = %v, want %v", tt.resBodyModel, tt.wantResBody)
				}
			}
		})
	}
}

func TestCashLaunchDeleteByID(t *testing.T) {
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
			name:         "ParamError",
			reqParam:     "x",
			resBodyModel: &model.Error{},
			wantResCode:  http.StatusBadRequest,
			wantResBody:  model.BadRequestParamValidate("Id invalid"),
		},
		{
			name:         "NotFoundError",
			reqParam:     "0",
			resBodyModel: &model.Error{},
			wantResCode:  http.StatusNotFound,
			wantResBody:  model.NotFound(controllerCashLaunchTitle),
		},
		{
			name:         "RepositoryError",
			reqParam:     "2",
			repoError:    true,
			resBodyModel: &model.Error{},
			wantResCode:  http.StatusInternalServerError,
			wantResBody:  model.InternalServerErrorRepositoryPersist(controllerCashLaunchTitle),
		},
		{
			name:        "Success",
			reqParam:    "2",
			wantResCode: http.StatusNoContent,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log := hclog.New(&hclog.LoggerOptions{Level: hclog.LevelFromString("OFF")})
			repository, _ := repository_in_memory.NewInMemory(tt.repoError)
			repositoryCashLaunch := repository.CashLaunch()
			usecaseCashLaunch := usecase.NewCashLaunch(repositoryCashLaunch)
			controllerCashLaunch := controller.NewCashLaunch(log, usecaseCashLaunch)

			url := fmt.Sprintf("/api/cash/launch/%v", tt.reqParam)

			req, _ := http.NewRequest(http.MethodDelete, url, nil)
			handler := http.HandlerFunc(controllerCashLaunch.DeleteByID)
			res := httptest.NewRecorder()

			handler.ServeHTTP(res, req)

			if !reflect.DeepEqual(res.Code, tt.wantResCode) {
				t.Errorf("DeleteByID() got res.code = %v, want %v", res.Code, tt.wantResCode)
			}

			if tt.resBodyModel != nil {
				json.NewDecoder(res.Body).Decode(&tt.resBodyModel)
			}

			if !reflect.DeepEqual(tt.resBodyModel, tt.wantResBody) {
				t.Errorf("DeleteByID() got res.body = %v, want %v", tt.resBodyModel, tt.wantResBody)
			}
		})
	}
}
