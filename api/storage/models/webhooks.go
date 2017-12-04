package models

import (
	"time"
)

// Webhook represents a Dairycart webhook
type Webhook struct {
	ID         uint64    `json:"id"`          // id
	URL        string    `json:"url"`         // url
	EventType  string    `json:"event_type"`  // event_type
	CreatedOn  time.Time `json:"created_on"`  // created_on
	UpdatedOn  NullTime  `json:"updated_on"`  // updated_on
	ArchivedOn NullTime  `json:"archived_on"` // archived_on
}
