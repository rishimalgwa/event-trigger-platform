package schemas

import (
	"time"

	"github.com/google/uuid"
)

type EventLogResponse struct {
	Events []EventDetail `json:"events"`
}

type EventDetail struct {
	ID         uuid.UUID `json:"id"`
	TriggerID  uuid.UUID `json:"trigger_id"`
	Type       string    `json:"type"`
	APIURL     *string   `json:"api_url,omitempty"`
	APIPayload any       `json:"api_payload,omitempty"`
	ExecutedAt time.Time `json:"executed_at"`
	IsTest     bool      `json:"is_test"`
}
