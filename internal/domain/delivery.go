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
	Payload                  []byte
	Status                   DeliveryStatus
	Attempts                 int
	LastError                string
	CreatedAt, UpdatedAt     time.Time
}
