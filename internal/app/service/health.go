package service

import (
	"context"

	"shortly/internal/app/repository"
)

// HealthChecker is an interface for health checks
type HealthChecker interface {
	Ping(ctx context.Context) error
}

// healthService is a service for health checks
type healthService struct {
	repo repository.Repository
}

// NewHealthService creates a new health service
func NewHealthService(repo repository.Repository) HealthChecker {
	return &healthService{
		repo: repo,
	}
}

// Ping checks the health of the database connection
func (s *healthService) Ping(ctx context.Context) error {
	if repo, ok := s.repo.(repository.HealthChecker); ok {
		return repo.Ping(ctx)
	}

	return nil
}
