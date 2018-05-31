package models

import (
	"time"
)

// WebhookExecutionLog represents a Dairycart webhook execution log
type WebhookExecutionLog struct {
	ID         uint64    `json:"id"`          // id
	WebhookID  uint64    `json:"webhook_id"`  // webhook_id
	StatusCode int       `json:"status_code"` // status_code
	Succeeded  bool      `json:"succeeded"`   // succeeded
	ExecutedOn time.Time `json:"executed_on"` // executed_on
}

// WebhookExecutionLogCreationInput is a struct to use for creating WebhookExecutionLogs
type WebhookExecutionLogCreationInput struct {
	StatusCode int  `json:"status_code,omitempty"` // status_code
	Succeeded  bool `json:"succeeded,omitempty"`   // succeeded
}

// WebhookExecutionLogUpdateInput is a struct to use for updating WebhookExecutionLogs
type WebhookExecutionLogUpdateInput struct {
	WebhookID  uint64 `json:"webhook_id,omitempty"`  // webhook_id
	StatusCode int    `json:"status_code,omitempty"` // status_code
	Succeeded  bool   `json:"succeeded,omitempty"`   // succeeded
}

type WebhookExecutionLogListResponse struct {
	ListResponse
	WebhookExecutionLogs []WebhookExecutionLog `json:"webhook_execution_logs"`
}
