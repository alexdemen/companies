package postgres

import (
	"context"
	"github.com/alexdemen/companies/app/stores"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Instance struct {
	pool *pgxpool.Pool
}

func (inst *Instance) GetExecutor() stores.Executor {
	return &Executor{}
}

func NewInstance(ctx context.Context) *Instance {
	pool, err := pgxpool.Connect(ctx, "postgres://postgres:12345678@127.0.0.1:5433/companies")
	if err != nil {

	}

	return &Instance{pool: pool}
}

type Executor struct {
}
