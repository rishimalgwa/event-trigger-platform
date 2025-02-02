package trigger

import (
	"github.com/google/uuid"
	"github.com/rishimalgwa/event-trigger-platform/api/schemas"
	"github.com/rishimalgwa/event-trigger-platform/pkg/models"
	"gorm.io/gorm"
)

type repo struct {
	DB *gorm.DB
}

// Save implements Repository.
func (r *repo) Save(trigger *schemas.CreateTriggerRequest) error {
	t := &models.Trigger{
		// UserID:              trigger.UserID,
		Type:                models.TriggerType(trigger.Type),
		ScheduleTime:        trigger.ScheduleTime,
		IntervalSecs:        trigger.IntervalSecs,
		IsRecurring:         *trigger.IsRecurring,
		APIURL:              trigger.APIURL,
		APIPayload:          trigger.APIPayload,
		NumberOfOccurrences: trigger.NumberOfOccurrences,
	}

	return r.DB.Save(t).Error

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

func NewPostgresRepo(db *gorm.DB) Repository {
	return &repo{
		DB: db,
	}
}
