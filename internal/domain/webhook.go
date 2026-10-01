package domain

import (
	"time"
)

type Webhook struct {
	id string
	URL string
	Events []string
	Secret string
	CreatedAt time.Time
}