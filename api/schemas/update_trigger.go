package schemas

import "github.com/go-playground/validator/v10"

type UpdateTriggerRequest struct {
	APIURL              *string `json:"api_url,omitempty" validate:"omitempty,url"`
	APIPayload          *string `json:"api_payload,omitempty" validate:"omitempty"`
	NumberOfOccurrences *int    `json:"number_of_occurrences,omitempty" validate:"omitempty,gte=1"`
	IntervalSecs        *int    `json:"interval_secs,omitempty" validate:"omitempty,gte=10"`

	ScheduleTime *string `json:"schedule_time,omitempty" validate:"required,rfc3339"`
}

func (u *UpdateTriggerRequest) Validate() []*ErrorResponse {
	validate := validator.New()
	validate.RegisterValidation("rfc3339", validateRFC3339)

	var errors []*ErrorResponse

	err := validate.Struct(u)
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
