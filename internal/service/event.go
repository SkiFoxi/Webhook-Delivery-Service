package service

import (
	"context"
	"log/slog"
	"time"
	"uuid"

	"github.com/SkiFoxi/Webhook-Delivery-Service/internal/domain"
)

type JobSubmitter interface {
	Submit(job domain.DeliveryJob) error
}

type EventService struct {
	webhookRepo  domain.WebhookRepository
	deliveryRepo domain.DeliveryRepository
	pool         JobSubmitter
}

func NewEventService(
	webhookRepo domain.WebhookRepository,
	deliveryRepo domain.DeliveryRepository,
	pool JobSubmitter,
) *EventService {
	return &EventService{
		webhookRepo:  webhookRepo,
		deliveryRepo: deliveryRepo,
		pool:         pool,
	}
}

func (s *EventService) Trigger(ctx context.Context, eventType string, payload []byte) error {
	webhooks, err := s.webhookRepo.ListByEvent(ctx, eventType)

	if err != nil {
		return err
	}

	if len(webhooks) == 0 {
		return nil
	}

	slog.Info("event triggered", "event_type", eventType, "webhooks_count", len(webhooks))

	for _, wh := range webhooks {
		d := &domain.Delivery{
			ID:        uuid.New().String(),
			WebhookID: wh.ID,
			EventType: eventType,
			Payload:   payload,
			Status:    domain.DelivPending,
			Attempts:  0,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if err := s.deliveryRepo.Create(ctx, d); err != nil {
			slog.Error("cannot create delivery", "err", err, "webhook_id", wh.ID)
			continue
		}
		job := domain.DeliveryJob{
			DeliveryID: d.ID,
			WebhookURL: wh.URL,
			Payload:    payload,
			Secret:     wh.Secret,
		}
		if err := s.pool.Submit(job); err != nil {
			slog.Error("cannot submit job", "err", err, "job_id", d.ID)
			continue
		}
	}
	return nil
}
