package service

import (
	"context"

	"shortly/internal/app/repository"
)

type HealthServiceInterface interface {
	Ping(ctx context.Context) error
}

type HealthService struct {
	repo repository.Repository
}

func NewHealthService(repo repository.Repository) *HealthService {
	return &HealthService{
		repo: repo,
	}
}

func (s *HealthService) Ping(ctx context.Context) error {
	return s.repo.Ping(ctx)
}
