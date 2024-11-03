package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgxPool interface {
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Close()
}

type DatabaseRepository struct {
	db PgxPool
}

func NewDatabaseRepository(ctx context.Context, dsn string) (*DatabaseRepository, error) {
	db, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}

	return &DatabaseRepository{db: db}, nil
}

func (d DatabaseRepository) Set(_ URL) error {
	return nil
}

func (d DatabaseRepository) Get(_ string) (*URL, bool) {
	return nil, false
}

func (d DatabaseRepository) CreateMemento() *Memento {
	var results []URL
	return &Memento{State: results}
}

func (d DatabaseRepository) Restore(_ *Memento) {
}

func (d DatabaseRepository) Ping(ctx context.Context) error {
	var result int
	return d.db.QueryRow(ctx, "SELECT 1").Scan(&result)
}
