package schemas

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

// created or updated
type CreateTriggerResponse struct {
	ID                  uuid.UUID  `json:"id"`
	UserID              uuid.UUID  `json:"user_id"`
	Type                string     `json:"type"`
	APIURL              *string    `json:"api_url,omitempty"`
	APIPayload          any        `json:"api_payload,omitempty"`
	Schedule            *time.Time `json:"schedule,omitempty"`
	NumberOfOccurrences *int       `json:"number_of_occurrences,omitempty"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
}

// CreateTriggerRequest struct with validation
type CreateTriggerRequest struct {
	Type                string  `json:"type" validate:"required,oneof=scheduled api"`
	APIURL              *string `json:"api_url,omitempty" validate:"omitempty,url"`
	APIPayload          *string `json:"api_payload,omitempty" validate:"omitempty,required_if=Type api"`
	NumberOfOccurrences *int    `json:"number_of_occurrences,omitempty" validate:"omitempty,gte=1"`
	IntervalSecs        *int    `json:"interval_secs,omitempty" validate:"omitempty,gte=10"`
	IsRecurring         *bool   `json:"is_recurring,omitempty"`
	ScheduleTime        *string `json:"schedule_time,omitempty" validate:"omitempty,rfc3339"`
}

// Custom Validator
func validateRFC3339(fl validator.FieldLevel) bool {
	if fl.Field().String() == "" {
		return true // Allow empty value
	}
	_, err := time.Parse(time.RFC3339, fl.Field().String())
	return err == nil
}

// Validate function
func (c *CreateTriggerRequest) Validate() []*ErrorResponse {
	validate := validator.New()
	validate.RegisterValidation("rfc3339", validateRFC3339) // Register custom validation

	var errors []*ErrorResponse

	err := validate.Struct(c)
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
