package route

import (
	"github.com/CharlesSchiavinato/minsait-challenge-backend/controller"
	"github.com/CharlesSchiavinato/minsait-challenge-backend/router"
	"github.com/CharlesSchiavinato/minsait-challenge-backend/service/database/repository"
	"github.com/CharlesSchiavinato/minsait-challenge-backend/usecase"
	"github.com/hashicorp/go-hclog"
)

type CashLaunchRouteParameters struct {
	AppRouter            router.Router
	Log                  hclog.Logger
	RepositoryCashLaunch repository.CashLaunch
}

func CashLaunchRoute(params *CashLaunchRouteParameters) {
	usecaseCashLaunch := usecase.NewCashLaunch(params.RepositoryCashLaunch)
	controllerCashLaunch := controller.NewCashLaunch(params.Log, usecaseCashLaunch)

	pathApiCashLaunch := "/api/cash/launch"
	pathApiCashLaunchParam := params.AppRouter.PathFormat("/api/cash/launch/%s", "param")

	params.AppRouter.Get(pathApiCashLaunch, controllerCashLaunch.List)
	params.AppRouter.Get(pathApiCashLaunchParam, controllerCashLaunch.GetByID)

	params.AppRouter.Post(pathApiCashLaunch, controllerCashLaunch.Insert)

	params.AppRouter.Put(pathApiCashLaunchParam, controllerCashLaunch.Update)

	params.AppRouter.Delete(pathApiCashLaunchParam, controllerCashLaunch.DeleteByID)
}
