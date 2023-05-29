package controller

import (
	"encoding/json"
	"net/http"

	"github.com/CharlesSchiavinato/minsait-challenge-backend/model"
)

type NotFound struct{}

func NewNotFound() *NotFound {
	return &NotFound{}
}

func (*NotFound) NotFound(rw http.ResponseWriter, req *http.Request) {
	modelNotFound := model.NotFound("URL")
	rw.WriteHeader(http.StatusNotFound)
	json.NewEncoder(rw).Encode(modelNotFound)
}
