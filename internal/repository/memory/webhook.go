package memory

import (
	"context"
	"slices"
	"sync"

	"github.com/SkiFoxi/Webhook-Delivery-Service/internal/domain"
)

type webhookRepo struct {
	mu       sync.RWMutex
	webhooks map[string]*domain.Webhook
}

func NewWebhookRepo() domain.WebhookRepository {
	webhooks := make(map[string]*domain.Webhook)
	return &webhookRepo{
		webhooks: webhooks,
	}
}

func (wr *webhookRepo) Create(ctx context.Context, w *domain.Webhook) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	wr.mu.Lock()
	defer wr.mu.Unlock()
	if _, ok := wr.webhooks[w.ID]; ok {
		return domain.ErrWebhookExists
	}
	copyW := *w
	wr.webhooks[w.ID] = &copyW
	return nil

}

func (wr *webhookRepo) GetByID(ctx context.Context, id string) (*domain.Webhook, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	wr.mu.RLock()
	defer wr.mu.RUnlock()
	arg, ok := wr.webhooks[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	copyArg := *arg
	return &copyArg, nil
}

func (wr *webhookRepo) ListByEvent(ctx context.Context, eventType string) ([]*domain.Webhook, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	result := make([]*domain.Webhook, 0)

	wr.mu.RLock()
	defer wr.mu.RUnlock()
	for _, arg := range wr.webhooks {
		if slices.Contains(arg.Events, eventType) {
			copyArg := *arg
			result = append(result, &copyArg)
		}
	}
	return result, nil
}

func (wr *webhookRepo) Delete(ctx context.Context, id string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	wr.mu.Lock()
	defer wr.mu.Unlock()
	if _, ok := wr.webhooks[id]; !ok {
		return domain.ErrNotFound
	}
	delete(wr.webhooks, id)
	return nil
}
