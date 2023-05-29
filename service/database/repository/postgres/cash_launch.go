package repository

import (
	"github.com/CharlesSchiavinato/minsait-challenge-backend/model"
	"github.com/CharlesSchiavinato/minsait-challenge-backend/service/database/repository"
	"github.com/lib/pq"
)

type PostgresCashLaunch struct {
	Postgres *Postgres
}

func NewCashLaunch(postgres *Postgres) repository.CashLaunch {
	return &PostgresCashLaunch{Postgres: postgres}
}

func (postgresCashLaunch *PostgresCashLaunch) Insert(modelCurrency *model.CashLaunch) (*model.CashLaunch, error) {
	query :=
		`INSERT INTO 
			cash_launch
			(reference_date, type, description, value, updated_at, created_at)
		VALUES
			($1, $2, $3, $4, $5, $6)
		RETURNING *;`

	row := postgresCashLaunch.Postgres.Conn.QueryRow(
		query,
		modelCurrency.ReferenceDate,
		modelCurrency.Type,
		modelCurrency.Description,
		modelCurrency.Value,
		modelCurrency.UpdatedAt,
		modelCurrency.CreatedAt,
	)

	modelCashLaunchInsert := &model.CashLaunch{}

	err := row.Scan(
		&modelCashLaunchInsert.ID,
		&modelCashLaunchInsert.ReferenceDate,
		&modelCashLaunchInsert.Type,
		&modelCashLaunchInsert.Description,
		&modelCashLaunchInsert.Value,
		&modelCashLaunchInsert.UpdatedAt,
		&modelCashLaunchInsert.CreatedAt,
	)

	// repository error duplicate key
	if errPQ, ok := err.(*pq.Error); ok {
		if errPQ.Code == "23505" {
			err = repository.ErrDuplicateKey{Message: errPQ.Detail}
		}
	}

	return modelCashLaunchInsert, err
}

func (postgresCashLaunch *PostgresCashLaunch) List() (model.CashLaunches, error) {
	query :=
		`SELECT
			id, reference_date, type, description, value, updated_at, created_at
		FROM
			cash_launch
		ORDER BY
			reference_date, type, value`

	rows, err := postgresCashLaunch.Postgres.Conn.Query(query)

	modelCashLaunches := model.CashLaunches{}

	if err != nil {
		return modelCashLaunches, err
	}

	defer rows.Close()

	for rows.Next() {
		modelCashLaunch := model.CashLaunch{}

		err = rows.Scan(
			&modelCashLaunch.ID,
			&modelCashLaunch.ReferenceDate,
			&modelCashLaunch.Type,
			&modelCashLaunch.Description,
			&modelCashLaunch.Value,
			&modelCashLaunch.UpdatedAt,
			&modelCashLaunch.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		modelCashLaunches = append(modelCashLaunches, modelCashLaunch)
	}

	return modelCashLaunches, err
}

func (postgresCashLaunch *PostgresCashLaunch) GetByID(id int64) (*model.CashLaunch, error) {
	query :=
		`SELECT
			id, reference_date, type, description, value, updated_at, created_at
		FROM
			cash_launch
		WHERE
			id = $1`

	row := postgresCashLaunch.Postgres.Conn.QueryRow(query, id)

	modelCashLaunch := model.CashLaunch{}

	err := row.Scan(
		&modelCashLaunch.ID,
		&modelCashLaunch.ReferenceDate,
		&modelCashLaunch.Type,
		&modelCashLaunch.Description,
		&modelCashLaunch.Value,
		&modelCashLaunch.UpdatedAt,
		&modelCashLaunch.CreatedAt,
	)

	// repository error not found
	if err != nil && err.Error() == "sql: no rows in result set" {
		err = repository.ErrNotFound{Message: err.Error()}
	}

	return &modelCashLaunch, err
}

func (postgresCashLaunch *PostgresCashLaunch) Update(modelCashLaunch *model.CashLaunch) (*model.CashLaunch, error) {
	query :=
		`UPDATE
		cash_launch
	SET
		reference_date = $2,
		type = $3,
		description = $4,
		value = $5,
		updated_at = $6
	WHERE
		id = $1
	RETURNING *;`

	row := postgresCashLaunch.Postgres.Conn.QueryRow(
		query,
		modelCashLaunch.ID,
		modelCashLaunch.ReferenceDate,
		modelCashLaunch.Type,
		modelCashLaunch.Description,
		modelCashLaunch.Value,
		modelCashLaunch.UpdatedAt,
	)

	modelCashLaunchUpdate := &model.CashLaunch{}

	err := row.Scan(
		&modelCashLaunchUpdate.ID,
		&modelCashLaunchUpdate.ReferenceDate,
		&modelCashLaunchUpdate.Type,
		&modelCashLaunchUpdate.Description,
		&modelCashLaunchUpdate.Value,
		&modelCashLaunchUpdate.UpdatedAt,
		&modelCashLaunchUpdate.CreatedAt,
	)

	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			// repository error not found
			err = repository.ErrNotFound{Message: err.Error()}
		} else if errPQ, ok := err.(*pq.Error); ok {
			// repository error duplicate key
			if errPQ.Code == "23505" {
				err = repository.ErrDuplicateKey{Message: errPQ.Detail}
			}
		}
	}

	return modelCashLaunchUpdate, err
}

func (postgresCashLaunch *PostgresCashLaunch) DeleteByID(id int64) error {
	query :=
		`DELETE FROM
		cash_launch
	WHERE
		id = $1`

	sqlResult, err := postgresCashLaunch.Postgres.Conn.Exec(query, id)

	if err == nil {
		rowsAffected, errRA := sqlResult.RowsAffected()

		if errRA != nil {
			err = errRA
		} else if rowsAffected == 0 {
			err = repository.ErrNotFound{}
		}
	}

	return err
}
