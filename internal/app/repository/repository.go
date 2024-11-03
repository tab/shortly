package repository

import (
	"context"

	"shortly/internal/logger"
)

type URL struct {
	UUID      string `json:"uuid"`
	LongURL   string `json:"long_url"`
	ShortCode string `json:"short_code"`
}

type Memento struct {
	State []URL `json:"state"`
}

type Repository interface {
	Set(url URL) error
	Get(shortCode string) (*URL, bool)
	CreateMemento() *Memento
	Restore(m *Memento)
	Ping(ctx context.Context) error
}

type Builder interface {
	CreateRepository(ctx context.Context) (Repository, error)
}

type Factory struct {
	DSN    string
	Logger *logger.Logger
}

func NewRepository(ctx context.Context, factory Builder) (Repository, error) {
	return factory.CreateRepository(ctx)
}

func (f *Factory) CreateRepository(ctx context.Context) (Repository, error) {
	db, err := NewDatabaseRepository(ctx, f.DSN)
	if err == nil {
		f.Logger.Info().Msg("Using PostgreSQL database")
		return db, nil
	}

	f.Logger.Info().Msg("Using in-memory repository")
	return NewInMemoryRepository(), nil
}
