package service

import (
	"context"
	"time"
	"uuid"

	"github.com/SkiFoxi/Webhook-Delivery-Service/internal/domain"
)

type WebhookService struct {
	repo domain.WebhookRepository
}

func NewWebhookService(repo domain.WebhookRepository) *WebhookService {
	return &WebhookService{
		repo: repo,
	}
}

func (s *WebhookService) Create(ctx context.Context, url string, events []string) (*domain.Webhook, error) {
	if url == "" {
		return nil, domain.ErrInvalidURL
	}

	if len(events) == 0 {
		return nil, domain.ErrNoEvents
	}

	w := &domain.Webhook{
		ID:        uuid.New().String(),
		URL:       url,
		Events:    events,
		Secret:    uuid.New().String(),
		CreatedAt: time.Now(),
	}
	if err := s.repo.Create(ctx, w); err != nil {
		return nil, err
	}

	return w, nil
}

func (s *WebhookService) GetByID(ctx context.Context, id string)(*domain.Webhook, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *WebhookService) ListByEvent(ctx context.Context, eventType string)([]*domain.Webhook, error) {
	return s.repo.ListByEvent(ctx, eventType)
}

func (s *WebhookService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
