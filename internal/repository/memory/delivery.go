package memory

import (
	"context"
	"sync"
	"time"

	"github.com/SkiFoxi/Webhook-Delivery-Service/internal/domain"
)

type deliveryRepo struct {
	mu         sync.RWMutex
	deliveries map[string]*domain.Delivery
}

func NewDeliveryRepo() domain.DeliveryRepository {
	deliveries := make(map[string]*domain.Delivery)
	return &deliveryRepo{
		deliveries: deliveries,
	}
}

func (dr *deliveryRepo) Create(ctx context.Context, d *domain.Delivery) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	dr.mu.Lock()
	defer dr.mu.Unlock()
	if _, ok := dr.deliveries[d.ID]; ok {
		return domain.ErrDeliveryExists
	}

	copyD := *d
	dr.deliveries[copyD.ID] = &copyD
	return nil
}

func (dr *deliveryRepo) UpdateStatus(ctx context.Context, id string, status domain.DeliveryStatus, lastError string) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	dr.mu.Lock()
	defer dr.mu.Unlock()
	arg, ok := dr.deliveries[id]
	if !ok {
		return domain.ErrNotFound
	}	
	arg.Status = status
	arg.LastError = lastError
	arg.UpdatedAt = time.Now()
	return nil

}

func (dr *deliveryRepo) ListByWebhook(ctx context.Context, webhookID string) ([]*domain.Delivery, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	result := make([]*domain.Delivery, 0)
	dr.mu.RLock()
	defer dr.mu.RUnlock()
	for _, arg := range dr.deliveries {
		if arg.WebhookID == webhookID {
			copyArg := *arg
			result = append(result, &copyArg)
		}
	}
	return result, nil
}
