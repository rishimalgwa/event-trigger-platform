package schemas

import (
	"time"

	"github.com/google/uuid"
)

// created or updated
type CreateTriggerResponse struct {
	ID         uuid.UUID `json:"id"`
	UserID     uuid.UUID `json:"user_id"`
	Type       string    `json:"type"`
	APIURL     *string   `json:"api_url,omitempty"`
	APIPayload any       `json:"api_payload,omitempty"`
	Schedule   *string   `json:"schedule,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
