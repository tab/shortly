package worker

import (
	"context"
	"sync"

	"shortly/internal/app/config"
	"shortly/internal/app/dto"
	"shortly/internal/app/repository"
	"shortly/internal/logger"
)

type Worker interface {
	Start(ctx context.Context)
	Stop()
	Add(req dto.BatchDeleteParams)
}

type worker struct {
	cfg    *config.Config
	repo   repository.Repository
	logger *logger.Logger
	queue  chan dto.BatchDeleteParams
	wg     sync.WaitGroup
}

func NewDeleteWorker(cfg *config.Config, repo repository.Repository, logger *logger.Logger) Worker {
	return &worker{
		cfg:    cfg,
		repo:   repo,
		logger: logger,
		queue:  make(chan dto.BatchDeleteParams),
	}
}

func (w *worker) Start(ctx context.Context) {
	w.logger.Info().Msgf("worker starting in %s", w.cfg.AppEnv)

	w.wg.Add(1)
	go w.run(ctx)
}

func (w *worker) Stop() {
	close(w.queue)
	w.wg.Wait()
}

func (w *worker) Add(req dto.BatchDeleteParams) {
	w.queue <- req
}

func (w *worker) run(ctx context.Context) {
	defer w.wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case req, ok := <-w.queue:
			if !ok {
				return
			}
			w.perform(ctx, req)
		}
	}
}

func (w *worker) perform(ctx context.Context, req dto.BatchDeleteParams) {
	shortCodes := unique(req.ShortCodes)
	err := w.repo.DeleteURLsByUserID(ctx, req.UserID, shortCodes)

	if err != nil {
		w.logger.Error().Err(err).Msgf("Error deleting URLs for user %s", req.UserID)
	} else {
		w.logger.Info().Msgf("Deleted URLs for user %s: %v", req.UserID, shortCodes)
	}
}

func unique(shortCodes []string) []string {
	codesMap := make(map[string]struct{})

	for _, code := range shortCodes {
		codesMap[code] = struct{}{}
	}

	uniqueCodes := make([]string, 0, len(codesMap))
	for code := range codesMap {
		uniqueCodes = append(uniqueCodes, code)
	}

	return uniqueCodes
}
