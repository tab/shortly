package service

import (
	"context"

	"shortly/internal/app/repository"
)

// StatsReporter is an interface for stats reporter
type StatsReporter interface {
	Counters(ctx context.Context) (int, int, error)
}

// statsService is a service for service metrics
type statsService struct {
	repo repository.Repository
}

// NewStatsReporter creates a new stats reporter instance
func NewStatsReporter(repo repository.Repository) StatsReporter {
	return &statsService{
		repo: repo,
	}
}

// Counters returns the number of URLs and users
func (m *statsService) Counters(ctx context.Context) (int, int, error) {
	if repo, ok := m.repo.(repository.StatsReporter); ok {
		return repo.Counters(ctx)
	}

	return 0, 0, nil
}
