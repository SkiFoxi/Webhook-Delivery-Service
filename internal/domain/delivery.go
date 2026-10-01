package domain

import "time"

type DeliveryStatus string

const (
	DelivPending DeliveryStatus = "panding"
	DelivSuccess DeliveryStatus = "success"
	DelivFailed  DeliveryStatus = "failed"
)

type Delivery struct {
	ID, WebhookId, EventType string
	Payload                  []byte
	Status                   DeliveryStatus
	Attempts                 int
	LastError                string
	CreatedAt, UpdatedAt     time.Time
}
