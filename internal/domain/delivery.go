package domain

import "time"

type DeliveryStatus string

const (
	DelivPending DeliveryStatus = "pending"
	DelivSuccess DeliveryStatus = "success"
	DelivFailed  DeliveryStatus = "failed"
)

type Delivery struct {
	ID, WebhookID, EventType string
	Payload                  []byte //Это JSON Данные тела события
	Status                   DeliveryStatus
	Attempts                 int //Сколько раз пытались отправить
	LastError                string
	CreatedAt, UpdatedAt     time.Time
}
