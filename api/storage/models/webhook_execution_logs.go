package models

import (
	"time"
)

// WebhookExecutionLog represents a Dairycart webhookexecutionlog
type WebhookExecutionLog struct {
	ID         uint64    `json:"id"`          // id
	WebhookID  uint64    `json:"webhook_id"`  // webhook_id
	StatusCode int       `json:"status_code"` // status_code
	Succeeded  bool      `json:"succeeded"`   // succeeded
	ExecutedOn time.Time `json:"executed_on"` // executed_on
}
