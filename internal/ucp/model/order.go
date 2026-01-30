package model

type OrderWebhookEvent struct {
	EventID   string            `json:"event_id"`
	EventType string            `json:"event_type"`
	Timestamp string            `json:"timestamp"`
	Order     OrderWebhookOrder `json:"order"`
}

type OrderWebhookOrder struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}
