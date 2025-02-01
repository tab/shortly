package worker

import (
	"context"
	"sync"

	"shortly/internal/app/config"
	"shortly/internal/app/dto"
	"shortly/internal/app/repository"
	"shortly/internal/logger"
)

// QueueSize is the size of the worker queue
const QueueSize = 100

// Worker is an interface for the delete worker
type Worker interface {
	Start()
	Stop()
	Add(req dto.BatchDeleteParams)
}

type worker struct {
	ctx    context.Context
	cfg    *config.Config
	repo   repository.Repository
	queue  chan dto.BatchDeleteParams
	logger *logger.Logger
	wg     sync.WaitGroup
}

// NewDeleteWorker creates a new delete worker instance
func NewDeleteWorker(ctx context.Context, cfg *config.Config, repo repository.Repository, logger *logger.Logger) Worker {
	queue := make(chan dto.BatchDeleteParams, QueueSize)

	return &worker{
		ctx:    ctx,
		cfg:    cfg,
		repo:   repo,
		queue:  queue,
		logger: logger,
	}
}

// Start starts the application worker
func (w *worker) Start() {
	w.logger.Info().Msgf("Worker starting in %s environment", w.cfg.AppEnv)

	w.wg.Add(1)
	go w.run()
}

// Stop stops the application worker
func (w *worker) Stop() {
	w.wg.Wait()
}

// Add adds a new request to the worker queue
func (w *worker) Add(req dto.BatchDeleteParams) {
	select {
	case <-w.ctx.Done():
		w.logger.Warn().Msg("Worker is stopped")
	default:
		w.queue <- req
	}
}

func (w *worker) run() {
	defer w.wg.Done()
	for {
		select {
		case <-w.ctx.Done():
			return
		case req := <-w.queue:
			w.perform(req)
		}
	}
}

func (w *worker) perform(req dto.BatchDeleteParams) {
	shortCodes := unique(req.ShortCodes)
	err := w.repo.DeleteURLsByUserID(w.ctx, req.UserID, shortCodes)

	if err != nil {
		w.logger.Error().Err(err).Msgf("Error deleting URLs for user %s", req.UserID)
	} else {
		w.logger.Info().Msgf("Deleted URLs for user %s: %v", req.UserID, shortCodes)
	}
}

// unique returns a unique slice of short codes
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
