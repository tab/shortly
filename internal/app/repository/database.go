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

type databaseRepo struct {
	queries *db.Queries
	db      *pgxpool.Pool
}

func NewDatabaseRepository(ctx context.Context, dsn string) (Database, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}

	queries := db.New(pool)

	return &databaseRepo{
		db:      pool,
		queries: queries,
	}, nil
}

func (d *databaseRepo) Set(ctx context.Context, url URL) error {
	_, err := d.queries.CreateURL(ctx, db.CreateURLParams{
		UUID:      url.UUID,
		LongURL:   url.LongURL,
		ShortCode: url.ShortCode,
	})

	return err
}

func (d *databaseRepo) Get(ctx context.Context, shortCode string) (*URL, bool) {
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

func (d *databaseRepo) Ping(ctx context.Context) error {
	_, err := d.queries.HealthCheck(ctx)
	return err
}

func (d *databaseRepo) Close() {
	d.db.Close()
}
