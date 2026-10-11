package worker_test

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"

	"github.com/SkiFoxi/Webhook-Delivery-Service/internal/domain"
	"github.com/SkiFoxi/Webhook-Delivery-Service/internal/worker"
)
type mockSender struct {
	mu sync.Mutex
	calls int
	err error
}

func (m *mockSender) Send(ctx context.Context, job domain.DeliveryJob) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.calls++
	return m.err
}

type mockDeliveryRepo struct {
    mu       sync.Mutex
    statuses map[string]domain.DeliveryStatus
}

func (m *mockDeliveryRepo) Create(ctx context.Context, d *domain.Delivery) error {
    return nil
}

func (m *mockDeliveryRepo) UpdateStatus(ctx context.Context, id string, status domain.DeliveryStatus, lastError string) error {
    m.mu.Lock()
    defer m.mu.Unlock()
    m.statuses[id] = status
    return nil
}

func (m *mockDeliveryRepo) ListByWebhook(ctx context.Context, webhookID string) ([]*domain.Delivery, error) {
    return nil, nil
}


func TestPool_ProcessesJobs(t *testing.T) {
    sender := &mockSender{}
    repo := &mockDeliveryRepo{statuses: make(map[string]domain.DeliveryStatus)}

    pool := worker.NewPool(3, 10, sender, repo)

    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    pool.Start(ctx)

    for i := 0; i < 10; i++ {
        job := domain.DeliveryJob{DeliveryID: fmt.Sprintf("d-%d", i)}
        if err := pool.Submit(job); err != nil {
            t.Fatalf("submit %d: %v", i, err)
        }
    }

    pool.Stop()

    if sender.calls != 10 {
        t.Errorf("calls = %d; want 10", sender.calls)
    }

    for i := 0; i < 10; i++ {
        id := fmt.Sprintf("d-%d", i)
        if repo.statuses[id] != domain.DelivSuccess {
            t.Errorf("delivery %s status = %s; want success", id, repo.statuses[id])
        }
    }
}

func TestPool_HandlesSendError(t *testing.T) {
    sender := &mockSender{err: errors.New("network error")}
    repo := &mockDeliveryRepo{statuses: make(map[string]domain.DeliveryStatus)}
    pool := worker.NewPool(2, 10, sender, repo)
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    pool.Start(ctx)
    err := pool.Submit(domain.DeliveryJob{DeliveryID: "d-1"})
	if err != nil {
		t.Fatalf("submit: %v", err)
	}
    pool.Stop()

    if repo.statuses["d-1"] != domain.DelivFailed {
        t.Errorf("status = %s; want failed", repo.statuses["d-1"])
    }
}

func TestPool_SubmitQueueFull(t *testing.T) {
    sender := &mockSender{}
    repo := &mockDeliveryRepo{}

    // queueSize = 1, но воркеры НЕ запущены
    pool := worker.NewPool(1, 1, sender, repo)

    if err := pool.Submit(domain.DeliveryJob{DeliveryID: "d-1"}); err != nil {
        t.Fatalf("first submit: %v", err)
    }

    // очередь должна быть полна
    err := pool.Submit(domain.DeliveryJob{DeliveryID: "d-2"})
    if err == nil {
        t.Error("expected error 'queue full', got nil")
    }
}

