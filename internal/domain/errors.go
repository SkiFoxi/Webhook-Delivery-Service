package domain

import(
	"errors"
)

var (
	ErrNotFound = errors.New("not found")
	ErrInvalidURL = errors.New("invalid url")
	ErrNoEvents = errors.New("at last one event required")
	ErrWebhookExists = errors.New("webhook already exists")
)

