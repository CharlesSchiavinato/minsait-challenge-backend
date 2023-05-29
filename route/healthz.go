package route

import (
	"github.com/CharlesSchiavinato/minsait-challenge-backend/controller"
	"github.com/CharlesSchiavinato/minsait-challenge-backend/router"
	"github.com/CharlesSchiavinato/minsait-challenge-backend/service/cache"
	"github.com/CharlesSchiavinato/minsait-challenge-backend/service/database/repository"
	"github.com/CharlesSchiavinato/minsait-challenge-backend/usecase"
	"github.com/hashicorp/go-hclog"
)

type HealthzRouteParameters struct {
	AppRouter  router.Router
	Log        hclog.Logger
	Repository repository.Repository
	Cache      cache.Cache
}

func HealthzRoute(params *HealthzRouteParameters) {
	useCaseHealthz := usecase.NewHealthz(params.Repository, params.Cache)
	controllerHealthz := controller.NewHealthz(params.Log, useCaseHealthz)

	params.AppRouter.Get("/api/healthz", controllerHealthz.Check)
}
