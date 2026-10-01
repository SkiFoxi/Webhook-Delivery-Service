package domain

import "context"

type WebhookRepository interface {
	Create(ctx context.Context, w *Webhook) error                          // Сохранить новый webhook
	GetByID(ctx context.Context, id string) (*Webhook, error)              // Найти по ID
	ListByEvent(ctx context.Context, eventType string) ([]*Webhook, error) //Список всех вебхуков и событий
	Delete(ctx context.Context, id string) error                           //Удаление по ID
}

type DeliveryRepository interface {
	Create(ctx context.Context, d *Delivery) error
	UpdateStatus(ctx context.Context, id string, status DeliveryStatus, lastError string) error
	ListByWebhook(ctx context.Context, webhookID string)([]*Delivery, error)
}
