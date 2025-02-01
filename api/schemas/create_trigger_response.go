package schemas

import "time"

// created or updated
type CreateTriggerResponse struct {
	ID         uint      `json:"id"`
	UserID     uint      `json:"user_id"`
	Type       string    `json:"type"`
	APIURL     *string   `json:"api_url,omitempty"`
	APIPayload any       `json:"api_payload,omitempty"`
	Schedule   *string   `json:"schedule,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
