package spec

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func TruncateTables(ctx context.Context, dsn string) error {
	err := RunQuery(ctx, dsn, "TRUNCATE TABLE urls RESTART IDENTITY CASCADE")
	if err != nil {
		return err
	}

	return nil
}

func RunQuery(ctx context.Context, dsn string, query string) error {
	db, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return err
	}

	_, err = db.Exec(ctx, query)
	if err != nil {
		return err
	}
	db.Close()

	return nil
}
