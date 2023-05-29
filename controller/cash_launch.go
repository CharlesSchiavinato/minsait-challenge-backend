package controller

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/CharlesSchiavinato/minsait-challenge-backend/model"
	"github.com/CharlesSchiavinato/minsait-challenge-backend/service/database/repository"
	logger "github.com/CharlesSchiavinato/minsait-challenge-backend/service/logger"
	"github.com/CharlesSchiavinato/minsait-challenge-backend/usecase"
	"github.com/hashicorp/go-hclog"
)

type CashLaunch struct {
	Title             string
	Log               hclog.Logger
	UseCaseCashLaunch usecase.CashLaunch
}

func NewCashLaunch(log hclog.Logger, useCaseCashLaunch usecase.CashLaunch) *CashLaunch {
	return &CashLaunch{
		Title:             "CashLaunch",
		Log:               log,
		UseCaseCashLaunch: useCaseCashLaunch,
	}
}

// Insert godoc
// @Summary      Adicionar
// @Description  Adiciona Lançamento
// @Tags         Lançamentos
// @Accept       json
// @Produce      json
// @Param        request   body      model.parametersCashLaunchWrapper  true  "Lançamento"
// @Success      201  {object}  model.CashLaunch
// @Failure      400  {object}  model.Error
// @Failure      500  {object}  model.Error
// @Router       /cash/launch [post]
func (controllerCashLaunch *CashLaunch) Insert(rw http.ResponseWriter, req *http.Request) {

	modelCashLaunch := &model.CashLaunch{}

	err := json.NewDecoder(req.Body).Decode(modelCashLaunch)

	if err != nil {
		responseError := model.BadRequestDeserialize(controllerCashLaunch.Title)

		logger.LogErrorRequest(controllerCashLaunch.Log, req, responseError.Message, err)

		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(responseError)
		return
	}

	modelCashLaunchInsert, err := controllerCashLaunch.UseCaseCashLaunch.Insert(modelCashLaunch)

	if err != nil {
		var responseError *model.Error

		if _, ok := err.(usecase.ErrModelValidate); ok {
			responseError = model.BadRequestModelValidate(controllerCashLaunch.Title, err.Error())

			rw.WriteHeader(http.StatusBadRequest)
		} else {
			responseError = model.InternalServerErrorRepositoryPersist(controllerCashLaunch.Title)

			logger.LogErrorRequest(controllerCashLaunch.Log, req, responseError.Message, err)

			rw.WriteHeader(http.StatusInternalServerError)
		}

		json.NewEncoder(rw).Encode(responseError)
		return
	}

	rw.WriteHeader(http.StatusCreated)
	json.NewEncoder(rw).Encode(modelCashLaunchInsert)
}

// List godoc
// @Summary      Listar
// @Description  Retorna uma lista de Lançamentos
// @Tags         Lançamentos
// @Accept       json
// @Produce      json
// @Success      200 {object}  model.CashLaunches
// @Failure      500  {object}  model.Error
// @Router       /cash/launch [get]
func (controllerCashLaunch *CashLaunch) List(rw http.ResponseWriter, req *http.Request) {
	modelCashLaunches, err := controllerCashLaunch.UseCaseCashLaunch.List()

	if err != nil {
		responseError := model.InternalServerErrorRepositoryLoad(controllerCashLaunch.Title)

		logger.LogErrorRequest(controllerCashLaunch.Log, req, responseError.Message, err)

		rw.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(rw).Encode(responseError)
		return
	}

	json.NewEncoder(rw).Encode(modelCashLaunches)
}

// GetByID godoc
// @Summary      Consultar
// @Description  Retorna um Lançamento
// @Tags         Lançamentos
// @Accept       json
// @Produce      json
// @Param        param   path      string  false  "Id do Lançamento" example("1")
// @Success      200 {object}  model.CashLaunch
// @Failure      400  {object}  model.Error
// @Failure      404  {object}  model.Error
// @Failure      500  {object}  model.Error
// @Router       /cash/launch/{id} [get]
func (controllerCashLaunch *CashLaunch) GetByID(rw http.ResponseWriter, req *http.Request) {
	param := strings.Split(req.URL.Path, "/")[4]

	id, err := strconv.ParseInt(param, 10, 64)

	if err != nil {
		responseError := model.BadRequestParamValidate("Id invalid")

		logger.LogErrorRequest(controllerCashLaunch.Log, req, responseError.Message, err)

		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(responseError)
		return
	}

	modelCashLaunch, err := controllerCashLaunch.UseCaseCashLaunch.GetByID(id)

	if err != nil {
		var responseError *model.Error

		if _, ok := err.(repository.ErrNotFound); ok {
			responseError = model.NotFound(controllerCashLaunch.Title)

			rw.WriteHeader(http.StatusNotFound)
		} else {
			responseError = model.InternalServerErrorRepositoryLoad(controllerCashLaunch.Title)

			logger.LogErrorRequest(controllerCashLaunch.Log, req, responseError.Message, err)

			rw.WriteHeader(http.StatusInternalServerError)
		}

		json.NewEncoder(rw).Encode(responseError)
		return
	}

	json.NewEncoder(rw).Encode(modelCashLaunch)
}

// Update godoc
// @Summary      Alterar
// @Description  Altera um Lançamento
// @Tags         Lançamentos
// @Accept       json
// @Produce      json
// @Param        param   path      string  false  "Id do Lançamento" example("1")
// @Param        request   body      model.parametersCashLaunchWrapper  true  "Lançamento"
// @Success      200 {object}  model.CashLaunch
// @Failure      400  {object}  model.Error
// @Failure      404  {object}  model.Error
// @Failure      500  {object}  model.Error
// @Router       /cash/launch/{id} [put]
func (controllerCashLaunch *CashLaunch) Update(rw http.ResponseWriter, req *http.Request) {
	param := strings.Split(req.URL.Path, "/")[4]

	id, err := strconv.ParseInt(param, 10, 64)

	if err != nil {
		responseError := model.BadRequestParamValidate("Id invalid")

		logger.LogErrorRequest(controllerCashLaunch.Log, req, responseError.Message, err)

		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(responseError)
		return
	}

	modelCashLaunch := &model.CashLaunch{}

	err = json.NewDecoder(req.Body).Decode(modelCashLaunch)

	if err != nil {
		responseError := model.BadRequestDeserialize(controllerCashLaunch.Title)

		logger.LogErrorRequest(controllerCashLaunch.Log, req, responseError.Message, err)

		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(responseError)
		return
	}

	modelCashLaunch.ID = id

	modelCashLaunchUpdate, err := controllerCashLaunch.UseCaseCashLaunch.Update(modelCashLaunch)

	if err != nil {
		var responseError *model.Error

		if _, ok := err.(usecase.ErrModelValidate); ok {
			responseError = model.BadRequestModelValidate(controllerCashLaunch.Title, err.Error())

			rw.WriteHeader(http.StatusBadRequest)
		} else if _, ok := err.(repository.ErrNotFound); ok {
			responseError = model.NotFound(controllerCashLaunch.Title)

			rw.WriteHeader(http.StatusNotFound)
		} else {
			responseError = model.InternalServerErrorRepositoryPersist(controllerCashLaunch.Title)

			logger.LogErrorRequest(controllerCashLaunch.Log, req, responseError.Message, err)

			rw.WriteHeader(http.StatusInternalServerError)
		}

		json.NewEncoder(rw).Encode(responseError)
		return
	}

	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(modelCashLaunchUpdate)
}

// DeleteByID godoc
// @Summary      Excluir
// @Description  Exclui um Lançamento
// @Tags         Lançamentos
// @Accept       json
// @Produce      json
// @Param        param   path      string  false  "Id do Lançamento" example("1")
// @Success      204
// @Failure      400  {object}  model.Error
// @Failure      404  {object}  model.Error
// @Failure      500  {object}  model.Error
// @Router       /cash/launch/{id} [delete]
func (controllerCashLaunch *CashLaunch) DeleteByID(rw http.ResponseWriter, req *http.Request) {
	param := strings.Split(req.URL.Path, "/")[4]

	id, err := strconv.ParseInt(param, 10, 64)

	if err != nil {
		responseError := model.BadRequestParamValidate("Id invalid")

		logger.LogErrorRequest(controllerCashLaunch.Log, req, responseError.Message, err)

		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(responseError)
		return
	}

	err = controllerCashLaunch.UseCaseCashLaunch.DeleteByID(id)

	if err != nil {
		var responseError *model.Error

		if _, ok := err.(repository.ErrNotFound); ok {
			responseError = model.NotFound(controllerCashLaunch.Title)

			rw.WriteHeader(http.StatusNotFound)
		} else {
			responseError = model.InternalServerErrorRepositoryPersist(controllerCashLaunch.Title)

			logger.LogErrorRequest(controllerCashLaunch.Log, req, responseError.Message, err)

			rw.WriteHeader(http.StatusInternalServerError)
		}

		json.NewEncoder(rw).Encode(responseError)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}
