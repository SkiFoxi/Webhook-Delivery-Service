package domain

import (
	"time"
)

type Webhook struct {
	Id string
	URL string
	Events []string
	Secret string
	CreatedAt time.Time
}