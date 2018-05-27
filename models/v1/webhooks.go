package models

import (
	"time"
)

// Webhook represents a Dairycart webhook
type Webhook struct {
	ID          uint64     `json:"id"`           // id
	URL         string     `json:"url"`          // url
	EventType   string     `json:"event_type"`   // event_type
	ContentType string     `json:"content_type"` // content_type
	CreatedOn   time.Time  `json:"created_on"`   // created_on
	UpdatedOn   *Dairytime `json:"updated_on"`   // updated_on
	ArchivedOn  *Dairytime `json:"archived_on"`  // archived_on
}

// WebhookCreationInput is a struct to use for creating Webhooks
type WebhookCreationInput struct {
	URL         string `json:"url,omitempty"`          // url
	EventType   string `json:"event_type,omitempty"`   // event_type
	ContentType string `json:"content_type,omitempty"` // content_type
}

// WebhookUpdateInput is a struct to use for updating Webhooks
type WebhookUpdateInput struct {
	URL         string `json:"url,omitempty"`          // url
	EventType   string `json:"event_type,omitempty"`   // event_type
	ContentType string `json:"content_type,omitempty"` // content_type
}

type WebhookListResponse struct {
	ListResponse
	Webhooks []Webhook `json:"webhooks"`
}
