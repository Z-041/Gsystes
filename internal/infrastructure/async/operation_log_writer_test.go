package async

import (
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gsystes/backend/internal/domain/entity"
	"github.com/gsystes/backend/internal/infrastructure/logger"
)

func init() {
	logger.InitForTesting()
}

type mockOpLogRepo struct {
	created  int32
	failMode bool
	mu       sync.Mutex
	entries  []*entity.OperationLog
}

func (m *mockOpLogRepo) Create(log *entity.OperationLog) error {
	if m.failMode {
		return errors.New("db error")
	}
	m.mu.Lock()
	m.entries = append(m.entries, log)
	m.mu.Unlock()
	atomic.AddInt32(&m.created, 1)
	return nil
}

func (m *mockOpLogRepo) FindByPage(page, pageSize int) ([]entity.OperationLog, int64, error) {
	return nil, 0, nil
}

func TestNewOperationLogWriter_Defaults(t *testing.T) {
	repo := &mockOpLogRepo{}
	w := NewOperationLogWriter(repo, 0, 0)
	if w.workers != 4 {
		t.Fatalf("expected 4 workers, got %d", w.workers)
	}
	if cap(w.queue) != 4096 {
		t.Fatalf("expected queue size 4096, got %d", cap(w.queue))
	}
}

func TestOperationLogWriter_Write(t *testing.T) {
	repo := &mockOpLogRepo{}
	w := NewOperationLogWriter(repo, 2, 64)
	w.Start()
	defer w.Stop()

	for i := 0; i < 50; i++ {
		w.Write(&entity.OperationLog{
			Method: "POST",
			Path:   "/test",
		})
	}

	time.Sleep(100 * time.Millisecond)

	count := atomic.LoadInt32(&repo.created)
	if count != 50 {
		t.Fatalf("expected 50 logs written, got %d", count)
	}
}

func TestOperationLogWriter_Write_Async(t *testing.T) {
	repo := &mockOpLogRepo{}
	w := NewOperationLogWriter(repo, 2, 64)
	w.Start()
	defer w.Stop()

	w.Write(&entity.OperationLog{Method: "GET", Path: "/async"})

	if atomic.LoadInt32(&repo.created) > 0 {
		t.Fatal("write should be async, log should not be persisted immediately")
	}

	time.Sleep(100 * time.Millisecond)

	if atomic.LoadInt32(&repo.created) != 1 {
		t.Fatal("log should be persisted after short wait")
	}
}

func TestOperationLogWriter_QueueFull_DropsEntry(t *testing.T) {
	repo := &mockOpLogRepo{failMode: true}
	w := NewOperationLogWriter(repo, 1, 2)
	w.Start()

	for i := 0; i < 10; i++ {
		w.Write(&entity.OperationLog{Method: "POST", Path: "/drop"})
	}

	time.Sleep(200 * time.Millisecond)
	w.Stop()

	count := atomic.LoadInt32(&repo.created)
	if count > 2 {
		t.Fatalf("queue should be bounded, expected <=2, got %d", count)
	}
}

func TestOperationLogWriter_Stop_DrainsQueue(t *testing.T) {
	repo := &mockOpLogRepo{}
	w := NewOperationLogWriter(repo, 1, 64)
	w.Start()

	for i := 0; i < 20; i++ {
		w.Write(&entity.OperationLog{Method: "PUT", Path: "/drain"})
	}

	w.Stop()

	count := atomic.LoadInt32(&repo.created)
	if count != 20 {
		t.Fatalf("Stop should drain remaining entries, expected 20, got %d", count)
	}
}
