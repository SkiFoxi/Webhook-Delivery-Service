package http

import "time"

type CreateWebhookRequest struct { //Должны принимать
	URL    string   `json:"url"`
	Events []string `json:"events"`
}

type WebhookResponse struct { //То, что будем отправлять
	ID        string    `json:"id"`
	URL       string    `json:"url"`
	Events    []string  `json:"events"`
	CreatedAt time.Time `json:"created_at"`
}
