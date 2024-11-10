package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"shortly/internal/app/repository/db"
)

type Database interface {
	Repository
	HealthChecker
	Close()
}

type DatabaseRepo struct {
	queries *db.Queries
	db      *pgxpool.Pool
}

func NewDatabaseRepository(ctx context.Context, dsn string) (Database, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}

	queries := db.New(pool)

	return &DatabaseRepo{
		db:      pool,
		queries: queries,
	}, nil
}

func (d *DatabaseRepo) Set(ctx context.Context, url URL) error {
	_, err := d.queries.CreateURL(ctx, db.CreateURLParams{
		UUID:      url.UUID,
		LongURL:   url.LongURL,
		ShortCode: url.ShortCode,
	})

	return err
}

func (d *DatabaseRepo) Get(ctx context.Context, shortCode string) (*URL, bool) {
	url, err := d.queries.GetURLByShortCode(ctx, shortCode)
	if err != nil {
		return nil, false
	}

	return &URL{
		UUID:      url.UUID,
		LongURL:   url.LongURL,
		ShortCode: url.ShortCode,
	}, true
}

func (d *DatabaseRepo) Ping(ctx context.Context) error {
	_, err := d.queries.HealthCheck(ctx)
	return err
}

func (d *DatabaseRepo) Close() {
	d.db.Close()
}
