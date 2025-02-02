package schemas

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

// created or updated
type CreateTriggerResponse struct {
	ID                  uuid.UUID `json:"id"`
	UserID              uuid.UUID `json:"user_id"`
	Type                string    `json:"type"`
	APIURL              *string   `json:"api_url,omitempty"`
	APIPayload          any       `json:"api_payload,omitempty"`
	Schedule            *string   `json:"schedule,omitempty"`
	NumberOfOccurrences *int      `json:"number_of_occurrences,omitempty"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

type CreateTriggerRequest struct {
	Type                string     `json:"type" validate:"required, oneof=scheduled api"`
	APIURL              *string    `json:"api_url,omitempty" validate:"url"`
	APIPayload          *string    `json:"api_payload,omitempty" validate:"required_if=Type api"`
	NumberOfOccurrences *int       `json:"number_of_occurrences,omitempty" validate:"gte=1"`
	IntervalSecs        *int       `json:"interval_secs,omitempty" validate:"gte=10"`
	IsRecurring         *bool      `json:"is_recurring,omitempty"`
	ScheduleTime        *time.Time `json:"schedule_time,omitempty" validate:"time"`
}

func (*CreateTriggerRequest) Validate() []*ErrorResponse {
	validate := validator.New()
	var errors []*ErrorResponse

	err := validate.Struct(CreateTriggerRequest{})
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			errors = append(errors, &ErrorResponse{
				FailedField: err.StructNamespace(),
				Tag:         err.Tag(),
				Value:       err.Param(),
			})
		}
	}
	return errors
}
