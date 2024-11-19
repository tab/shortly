package repository

import (
	"context"

	"github.com/google/uuid"

	"shortly/internal/logger"
)

type URL struct {
	UUID      uuid.UUID `json:"uuid"`
	LongURL   string    `json:"long_url"`
	ShortCode string    `json:"short_code"`
	UserUUID  uuid.UUID `json:"user_uuid"`
}

type User struct {
	UUID uuid.UUID `json:"uuid"`
}

type Memento struct {
	State []URL `json:"state"`
}

type Repository interface {
	CreateURL(ctx context.Context, url URL) (*URL, error)
	CreateURLs(ctx context.Context, urls []URL) error
	GetURLByShortCode(ctx context.Context, shortCode string) (*URL, bool)
	GetURLsByUserID(ctx context.Context, uuid uuid.UUID, limit, offset int64) ([]URL, int, error)
}

type HealthChecker interface {
	Ping(ctx context.Context) error
}

type Factory struct {
	DSN    string
	Logger *logger.Logger
}

func NewRepository(ctx context.Context, f *Factory) (Repository, error) {
	if f.DSN != "" {
		db, err := NewDatabaseRepository(ctx, f.DSN)
		if err == nil {
			f.Logger.Info().Msg("Using PostgreSQL database")
			return db, nil
		}
	}

	f.Logger.Info().Msg("Using in-memory repository")
	return NewInMemoryRepository(), nil
}
