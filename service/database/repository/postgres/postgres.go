package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/CharlesSchiavinato/minsait-challenge-backend/service/database/repository"
	"github.com/CharlesSchiavinato/minsait-challenge-backend/util"
)

type Postgres struct {
	Conn *sql.DB
}

func NewPostgres(config *util.Config) (repository.Repository, error) {
	db, err := sql.Open(config.DBDriver, config.DBURL)

	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(5)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err == nil {
		err = db.PingContext(ctx)
	}

	postgres := &Postgres{Conn: db}

	NewCashLaunch(postgres)

	return postgres, err
}

func (postgres *Postgres) Check() error {
	return postgres.Conn.Ping()
}

func (postgres *Postgres) Close() error {
	return postgres.Conn.Close()
}

func (postgres *Postgres) CashLaunch() repository.CashLaunch {
	return NewCashLaunch(postgres)
}

func (postgres *Postgres) CashBalanceDaily() repository.CashBalanceDaily {
	return NewCashBalanceDaily(postgres)
}
