package controller

import (
	"encoding/json"
	"net/http"

	"github.com/CharlesSchiavinato/minsait-challenge-backend/model"
	logger "github.com/CharlesSchiavinato/minsait-challenge-backend/service/logger"
	"github.com/CharlesSchiavinato/minsait-challenge-backend/usecase"
	"github.com/hashicorp/go-hclog"
)

type Healthz struct {
	Title          string
	Log            hclog.Logger
	UseCaseHealthz usecase.Healthz
}

func NewHealthz(log hclog.Logger, useCaseHealthz usecase.Healthz) *Healthz {
	return &Healthz{
		Title:          "Health Check",
		Log:            log,
		UseCaseHealthz: useCaseHealthz,
	}
}

func (controllerHealthz *Healthz) Check(rw http.ResponseWriter, req *http.Request) {
	err := controllerHealthz.UseCaseHealthz.CheckRepository()

	modelHealthz := &model.Healthz{}

	if err != nil {
		logger.LogErrorRequest(controllerHealthz.Log, req, controllerHealthz.Title, err)
		rw.WriteHeader(http.StatusFailedDependency)
		modelHealthz.Database = err.Error()
	} else {
		modelHealthz.Database = "OK"
	}

	err = controllerHealthz.UseCaseHealthz.CheckCache()

	if err != nil {
		logger.LogErrorRequest(controllerHealthz.Log, req, controllerHealthz.Title, err)
		rw.WriteHeader(http.StatusFailedDependency)
		modelHealthz.Cache = err.Error()
	} else {
		modelHealthz.Cache = "OK"
	}

	json.NewEncoder(rw).Encode(modelHealthz)
}
