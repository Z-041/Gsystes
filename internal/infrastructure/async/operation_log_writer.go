package async

import (
	"sync"

	"github.com/gsystes/backend/internal/domain/entity"
	"github.com/gsystes/backend/internal/domain/repository"
	"github.com/gsystes/backend/internal/infrastructure/logger"
)

type OperationLogWriter struct {
	repo    repository.OperationLogRepository
	queue   chan *entity.OperationLog
	workers int
	wg      sync.WaitGroup
	stopCh  chan struct{}
}

func NewOperationLogWriter(repo repository.OperationLogRepository, workers int, queueSize int) *OperationLogWriter {
	if workers <= 0 {
		workers = 4
	}
	if queueSize <= 0 {
		queueSize = 4096
	}
	return &OperationLogWriter{
		repo:    repo,
		queue:   make(chan *entity.OperationLog, queueSize),
		workers: workers,
		stopCh:  make(chan struct{}),
	}
}

func (w *OperationLogWriter) Start() {
	for i := 0; i < w.workers; i++ {
		w.wg.Add(1)
		go w.worker(i)
	}
	logger.Info("async log writer started",
		logger.IntField("workers", w.workers),
		logger.IntField("queue_size", cap(w.queue)),
	)
}

func (w *OperationLogWriter) worker(id int) {
	defer w.wg.Done()
	for {
		select {
		case entry := <-w.queue:
			w.safeCreate(entry)
		case <-w.stopCh:
			w.drainRemaining()
			return
		}
	}
}

func (w *OperationLogWriter) safeCreate(entry *entity.OperationLog) {
	defer func() {
		if r := recover(); r != nil {
			logger.Error("log writer panic recovered",
				logger.AnyField("panic", r),
			)
		}
	}()
	if err := w.repo.Create(entry); err != nil {
		logger.Error("failed to write operation log", logger.ErrorField(err))
	}
}

func (w *OperationLogWriter) drainRemaining() {
	for {
		select {
		case entry := <-w.queue:
			w.safeCreate(entry)
		default:
			return
		}
	}
}

func (w *OperationLogWriter) Write(entry *entity.OperationLog) {
	select {
	case w.queue <- entry:
	default:
		logger.Warn("operation log queue full, dropping entry")
	}
}

func (w *OperationLogWriter) Stop() {
	close(w.stopCh)
	w.wg.Wait()
	logger.Info("async log writer stopped")
}
