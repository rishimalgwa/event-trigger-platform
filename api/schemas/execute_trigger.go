package schemas

import "time"

type ExecuteTriggerResponse struct {
	TriggerID  uint      `json:"trigger_id"`
	Type       string    `json:"type"`
	Status     string    `json:"status"` // "success" or "failed"
	ExecutedAt time.Time `json:"executed_at"`
	Error      *string   `json:"error,omitempty"`
}
