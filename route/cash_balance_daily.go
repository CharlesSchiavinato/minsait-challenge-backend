package route

import (
	"github.com/CharlesSchiavinato/minsait-challenge-backend/controller"
	"github.com/CharlesSchiavinato/minsait-challenge-backend/router"
	"github.com/CharlesSchiavinato/minsait-challenge-backend/service/database/repository"
	"github.com/CharlesSchiavinato/minsait-challenge-backend/usecase"
	"github.com/hashicorp/go-hclog"
)

type CashBalanceDailyRouteParameters struct {
	AppRouter                  router.Router
	Log                        hclog.Logger
	RepositoryCashBalanceDaily repository.CashBalanceDaily
}

func CashBalanceDailyRoute(params *CashBalanceDailyRouteParameters) {
	usecaseCashBalanceDaily := usecase.NewCashBalanceDaily(params.RepositoryCashBalanceDaily)
	controllerCashBalanceDaily := controller.NewCashBalanceDaily(params.Log, usecaseCashBalanceDaily)

	pathApiCashBalanceDaily := "/api/cash/balance/daily"
	pathApiCashBalanceDailyParam := params.AppRouter.PathFormat("/api/cash/balance/daily/%s", "param")

	params.AppRouter.Get(pathApiCashBalanceDaily, controllerCashBalanceDaily.GetByRangeReferenceDate)
	params.AppRouter.Get(pathApiCashBalanceDailyParam, controllerCashBalanceDaily.GetByReferenceDate)
}
