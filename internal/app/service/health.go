package service

import (
	"context"

	"shortly/internal/app/repository"
)

type HealthChecker interface {
	Ping(ctx context.Context) error
}

type healthService struct {
	repo repository.Repository
}

func NewHealthService(repo repository.Repository) HealthChecker {
	return &healthService{
		repo: repo,
	}
}

func (s *healthService) Ping(ctx context.Context) error {
	if repo, ok := s.repo.(repository.HealthChecker); ok {
		return repo.Ping(ctx)
	}

	return nil
}
