package controller

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/CharlesSchiavinato/minsait-challenge-backend/model"
	logger "github.com/CharlesSchiavinato/minsait-challenge-backend/service/logger"
	"github.com/CharlesSchiavinato/minsait-challenge-backend/usecase"
	"github.com/hashicorp/go-hclog"
)

type CashBalanceDaily struct {
	Title                   string
	Log                     hclog.Logger
	UseCaseCashBalanceDaily usecase.CashBalanceDaily
}

func NewCashBalanceDaily(log hclog.Logger, useCaseCashBalanceDaily usecase.CashBalanceDaily) *CashBalanceDaily {
	return &CashBalanceDaily{
		Title:                   "CashBalanceDaily",
		Log:                     log,
		UseCaseCashBalanceDaily: useCaseCashBalanceDaily,
	}
}

// GetByReferenceDate godoc
// @Summary      Consultar
// @Description  Retorna o Saldo de todos os Lançamentos realizado na Data Informada. Se não for encontrado nenhum lançamento então será retornado com saldo zero.
// @Tags         Saldo Diário
// @Accept       json
// @Produce      json
// @Param        date   path      string  false  "Data de Referencia (AAAA-MM-DD)" example("2020-05-23")
// @Success      200  {object}  model.CashBalanceDaily
// @Failure      400  {object}  model.Error
// @Failure      500  {object}  model.Error
// @Router       /cash/balance/daily/{date} [get]
func (controllerCashBalanceDaily *CashBalanceDaily) GetByReferenceDate(rw http.ResponseWriter, req *http.Request) {
	params := strings.Split(req.URL.Path, "/")

	referenceDateParam := ""

	if len(params) > 5 {
		referenceDateParam = params[5]
	}

	referenceDate, err := time.Parse("2006-01-02", referenceDateParam)

	if err != nil {
		responseError := model.BadRequestParamValidate("Date invalid")

		logger.LogErrorRequest(controllerCashBalanceDaily.Log, req, responseError.Message, err)

		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(responseError)
		return
	}

	modelCashBalanceDaily, err := controllerCashBalanceDaily.UseCaseCashBalanceDaily.GetByReferenceDate(referenceDate)

	if err != nil {
		var responseError *model.Error

		if _, ok := err.(usecase.ErrParamValidate); ok {
			responseError = model.BadRequestParamValidate(err.Error())

			rw.WriteHeader(http.StatusBadRequest)
		} else {
			responseError = model.InternalServerErrorRepositoryLoad(controllerCashBalanceDaily.Title)

			logger.LogErrorRequest(controllerCashBalanceDaily.Log, req, responseError.Message, err)

			rw.WriteHeader(http.StatusInternalServerError)
		}

		json.NewEncoder(rw).Encode(responseError)
		return
	}

	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(modelCashBalanceDaily)
}

// GetByRangeReferenceDate godoc
// @Summary      Consultar por Período
// @Description  Retorna o Saldo Diário de todos os Lançamentos realizado no Período informado. O período não pode ser superior a 31 dias.
// @Tags         Saldo Diário
// @Accept       json
// @Produce      json
// @Param        from query      string  true  "Data de Referencia Inicial (AAAA-MM-DD)" example("2020-05-23")
// @Param        to   query      string  true  "Data de Referencia Final (AAAA-MM-DD)" example("2020-05-23")
// @Success      200  {object}  model.CashBalanceDailies
// @Failure      400  {object}  model.Error
// @Failure      500  {object}  model.Error
// @Router       /cash/balance/daily [get]
func (controllerCashBalanceDaily *CashBalanceDaily) GetByRangeReferenceDate(rw http.ResponseWriter, req *http.Request) {
	cashBalanceDailyRangeReferenceDate, err := extractURLQueryParamsRangeReferenceDate(req)

	if err != nil {
		responseError := model.BadRequestParamValidate(err.Error())

		logger.LogErrorRequest(controllerCashBalanceDaily.Log, req, responseError.Message, err)

		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(responseError)
		return
	}

	modelCashBalanceDailies, err := controllerCashBalanceDaily.UseCaseCashBalanceDaily.GetByRangeReferenceDate(cashBalanceDailyRangeReferenceDate)

	if err != nil {
		var responseError *model.Error

		if _, ok := err.(usecase.ErrParamValidate); ok {
			responseError = model.BadRequestParamValidate(err.Error())

			rw.WriteHeader(http.StatusBadRequest)
		} else {
			responseError = model.InternalServerErrorRepositoryLoad(controllerCashBalanceDaily.Title)

			logger.LogErrorRequest(controllerCashBalanceDaily.Log, req, responseError.Message, err)

			rw.WriteHeader(http.StatusInternalServerError)
		}

		json.NewEncoder(rw).Encode(responseError)
		return
	}

	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(modelCashBalanceDailies)

}

func extractURLQueryParamsRangeReferenceDate(req *http.Request) (*model.CashBalanceDailyRangeReferenceDate, error) {
	fromParam := req.URL.Query().Get("from")
	toParam := req.URL.Query().Get("to")

	cashBalanceDailyRangeReferenceDate := &model.CashBalanceDailyRangeReferenceDate{}

	messages := []string{}

	if fromParam == "" {
		messages = append(messages, "The param from is empty")
	} else {
		referenceDateFrom, err := time.Parse("2006-01-02", fromParam)

		if err != nil {
			messages = append(messages, "The param from is invalid")
		}

		cashBalanceDailyRangeReferenceDate.From = referenceDateFrom
	}

	if toParam == "" {
		messages = append(messages, "The param to is empty")
	} else {
		referenceDateTo, err := time.Parse("2006-01-02", toParam)

		if err != nil {
			messages = append(messages, "The param to is invalid")
		}

		cashBalanceDailyRangeReferenceDate.To = referenceDateTo
	}

	if len(messages) > 0 {
		return nil, errors.New(strings.Join(messages, ";"))
	}

	return cashBalanceDailyRangeReferenceDate, nil
}
