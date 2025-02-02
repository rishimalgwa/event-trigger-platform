package trigger

import (
	"time"

	"github.com/google/uuid"
	"github.com/rishimalgwa/event-trigger-platform/api/schemas"
	"github.com/rishimalgwa/event-trigger-platform/pkg/models"
	"gorm.io/gorm"
)

type repo struct {
	DB *gorm.DB
}

func NewPostgresRepo(db *gorm.DB) Repository {
	return &repo{
		DB: db,
	}
}

// Save implements Repository.
func (r *repo) Save(trigger *schemas.CreateTriggerRequest) (*uuid.UUID, error) {
	triggerID := uuid.New()
	var schTime time.Time
	if trigger.ScheduleTime != nil {
		schTime, _ = time.Parse(time.RFC3339, *trigger.ScheduleTime)
	}
	t := &models.Trigger{
		// UserID:              trigger.UserID,
		BaseModel: models.BaseModel{
			ID: triggerID,
		},
		Type:                models.TriggerType(trigger.Type),
		ScheduleTime:        &schTime,
		IntervalSecs:        trigger.IntervalSecs,
		IsRecurring:         trigger.IsRecurring,
		APIURL:              trigger.APIURL,
		APIPayload:          trigger.APIPayload,
		NumberOfOccurrences: trigger.NumberOfOccurrences,
	}
	err := r.DB.Create(t).Error
	if err != nil {
		return nil, err
	}
	return &triggerID, nil

}

func (r *repo) Find(id *uuid.UUID) (*models.Trigger, error) {
	u := &models.Trigger{}
	result := r.DB.Where("id = ?", id).First(u)
	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return u, nil
}
