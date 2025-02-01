package schemas

import (
	"time"

	"github.com/google/uuid"
)

// all active triggers
type ListTriggersResponse struct {
	Triggers []TriggerDetail `json:"triggers"`
}

type TriggerDetail struct {
	ID        uuid.UUID `json:"id"`
	Type      string    `json:"type"`
	APIURL    *string   `json:"api_url,omitempty"`
	Schedule  *string   `json:"schedule,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
