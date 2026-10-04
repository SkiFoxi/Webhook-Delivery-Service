package domain

import (
	"time"
)

type Webhook struct {
	ID string
	URL string
	Events []string
	Secret string
	CreatedAt time.Time
}