package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"shortly/internal/app/repository/db"
)

// Database is an interface for database operations
type Database interface {
	Repository
	HealthChecker
	Close()
}

// DatabaseRepo is a repository for database operations
type DatabaseRepo struct {
	queries *db.Queries
	db      *pgxpool.Pool
}

// NewDatabaseRepository creates a new database repository instance
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

// CreateURL creates a new URL record
func (d *DatabaseRepo) CreateURL(ctx context.Context, url URL) (*URL, error) {
	row, err := d.queries.CreateURL(ctx, db.CreateURLParams{
		UUID:      url.UUID,
		LongURL:   url.LongURL,
		ShortCode: url.ShortCode,
		UserUUID:  url.UserUUID,
	})

	if err != nil {
		return nil, err
	}

	return &URL{
		UUID:      row.UUID,
		LongURL:   row.LongURL,
		ShortCode: row.ShortCode,
	}, nil
}

// CreateURLs creates new URL records
func (d *DatabaseRepo) CreateURLs(ctx context.Context, urls []URL) error {
	tx, err := d.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	q := d.queries.WithTx(tx)

	for _, url := range urls {
		_, err := q.CreateURL(ctx, db.CreateURLParams{
			UUID:      url.UUID,
			LongURL:   url.LongURL,
			ShortCode: url.ShortCode,
		})
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

// GetURLByShortCode returns a URL record by short code
func (d *DatabaseRepo) GetURLByShortCode(ctx context.Context, shortCode string) (*URL, bool) {
	row, err := d.queries.GetURLByShortCode(ctx, shortCode)
	if err != nil {
		return nil, false
	}

	return &URL{
		UUID:      row.UUID,
		LongURL:   row.LongURL,
		ShortCode: row.ShortCode,
		DeletedAt: row.DeletedAt.Time,
	}, true
}

// GetURLByUUID returns a URL record by UUID
func (d *DatabaseRepo) GetURLsByUserID(ctx context.Context, id uuid.UUID, limit, offset int64) ([]URL, int, error) {
	params := db.GetURLsByUserIDParams{
		UserUUID: id,
		Limit:    limit,
		Offset:   offset,
	}

	rows, err := d.queries.GetURLsByUserID(ctx, params)
	if err != nil {
		return nil, 0, err
	}

	urls := make([]URL, 0, len(rows))
	var total int

	if len(rows) > 0 {
		total = int(rows[0].Total)
	}

	for _, row := range rows {
		urls = append(urls, URL{
			UUID:      row.UUID,
			LongURL:   row.LongURL,
			ShortCode: row.ShortCode,
		})
	}

	return urls, total, nil
}

// DeleteURLsByUserID deletes URL records by user ID
func (d *DatabaseRepo) DeleteURLsByUserID(ctx context.Context, id uuid.UUID, shortCodes []string) error {
	return d.queries.DeleteURLsByUserIDAndShortCodes(ctx, db.DeleteURLsByUserIDAndShortCodesParams{
		UserUUID:   id,
		ShortCodes: shortCodes,
	})
}

// Ping checks the database connection
func (d *DatabaseRepo) Ping(ctx context.Context) error {
	_, err := d.queries.HealthCheck(ctx)
	return err
}

// Close closes the database connection
func (d *DatabaseRepo) Close() {
	d.db.Close()
}
