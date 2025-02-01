package schemas

import "time"

type EventLogResponse struct {
	Events []EventDetail `json:"events"`
}

type EventDetail struct {
	ID         uint      `json:"id"`
	TriggerID  uint      `json:"trigger_id"`
	Type       string    `json:"type"`
	APIURL     *string   `json:"api_url,omitempty"`
	APIPayload any       `json:"api_payload,omitempty"`
	ExecutedAt time.Time `json:"executed_at"`
	IsTest     bool      `json:"is_test"`
}
