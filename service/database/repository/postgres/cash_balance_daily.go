package repository

import (
	"time"

	"github.com/CharlesSchiavinato/minsait-challenge-backend/model"
	"github.com/CharlesSchiavinato/minsait-challenge-backend/service/database/repository"
	"github.com/CharlesSchiavinato/minsait-challenge-backend/util"
)

type PostgresCashBalanceDaily struct {
	Postgres *Postgres
}

func NewCashBalanceDaily(postgres *Postgres) repository.CashBalanceDaily {
	return &PostgresCashBalanceDaily{Postgres: postgres}
}

func (postgresCashBalanceDaily *PostgresCashBalanceDaily) GetByReferenceDate(referenceDate time.Time) (*model.CashBalanceDaily, error) {
	query :=
		`SELECT 
			reference_date, 
			SUM(CASE WHEN type = 'C' THEN value ELSE (value * -1) END) AS value
		FROM 
			cash_launch
		WHERE
			reference_date = $1
		GROUP BY 
			reference_date `

	row := postgresCashBalanceDaily.Postgres.Conn.QueryRow(query, referenceDate)

	modelCashBalance := model.CashBalanceDaily{}

	err := row.Scan(
		&modelCashBalance.ReferenceDate,
		&modelCashBalance.Value,
	)

	// repository error not found
	if err != nil && err.Error() == "sql: no rows in result set" {
		err = repository.ErrNotFound{Message: err.Error()}
	}

	return &modelCashBalance, err
}

func (postgresCashBalanceDaily *PostgresCashBalanceDaily) GetByRangeReferenceDate(cashBalanceGetByRangeReferenceDateParams *model.CashBalanceDailyRangeReferenceDate) (model.CashBalanceDailies, error) {
	query :=
		`SELECT 
			reference_date, 
			SUM(CASE WHEN type = 'C' THEN value ELSE (value * -1) END) AS value
		FROM 
			cash_launch
		WHERE
			reference_date BETWEEN $1 AND $2
		GROUP BY 
			reference_date `

	rows, err := postgresCashBalanceDaily.Postgres.Conn.Query(query, cashBalanceGetByRangeReferenceDateParams.From, cashBalanceGetByRangeReferenceDateParams.To)

	modelCashBalances := model.CashBalanceDailies{}

	if err != nil {
		return modelCashBalances, err
	}

	defer rows.Close()

	for rows.Next() {
		modelCashBalance := model.CashBalanceDaily{}

		err = rows.Scan(
			&modelCashBalance.ReferenceDate,
			&modelCashBalance.Value,
		)

		if err != nil {
			return nil, err
		}

		modelCashBalance.Value = util.MathRoundPrecision(modelCashBalance.Value, 2)

		modelCashBalances = append(modelCashBalances, modelCashBalance)
	}

	return modelCashBalances, err
}
