package repository

import (
	"context"
	"time"

	"github.com/google/uuid"

	"shortly/internal/logger"
)

// URL is a URL entity
type URL struct {
	UUID      uuid.UUID `json:"uuid"`
	LongURL   string    `json:"long_url"`
	ShortCode string    `json:"short_code"`
	UserUUID  uuid.UUID `json:"user_uuid"`
	DeletedAt time.Time `json:"deleted_at"`
}

// User is a user entity
type User struct {
	UUID uuid.UUID `json:"uuid"`
}

// Memento is a memento entity
type Memento struct {
	State []URL `json:"state"`
}

// Repository is an interface for repository
type Repository interface {
	CreateURL(ctx context.Context, url URL) (*URL, error)
	CreateURLs(ctx context.Context, urls []URL) error
	GetURLByShortCode(ctx context.Context, shortCode string) (*URL, bool)
	GetURLsByUserID(ctx context.Context, uuid uuid.UUID, limit, offset int64) ([]URL, int, error)
	DeleteURLsByUserID(ctx context.Context, uuid uuid.UUID, shortCodes []string) error
}

// HealthChecker is an interface for health checker
type HealthChecker interface {
	Ping(ctx context.Context) error
}

// StatsReporter is an interface for stats reporter
type StatsReporter interface {
	Counters(ctx context.Context) (int, int, error)
}

// Factory is a factory for repository
type Factory struct {
	DSN    string
	Logger *logger.Logger
}

// NewRepository creates a new repository instance
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
