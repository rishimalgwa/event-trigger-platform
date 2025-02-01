package schemas

import (
	"time"

	"github.com/google/uuid"
)

type ExecuteTriggerResponse struct {
	TriggerID  uuid.UUID `json:"trigger_id"`
	Type       string    `json:"type"`
	Status     string    `json:"status"` // "success" or "failed"
	ExecutedAt time.Time `json:"executed_at"`
	Error      *string   `json:"error,omitempty"`
}
