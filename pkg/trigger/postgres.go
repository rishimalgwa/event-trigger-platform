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

func (r *repo) FindAll() ([]*models.Trigger, error) {
	var triggers []*models.Trigger

	// get only triggers with execution status "initialized"
	return triggers, r.DB.Where("execution_status = ?", models.Initialized).Find(&triggers).Error
}

// Delete implements Repository.
func (r *repo) Delete(id *uuid.UUID) error {
	return r.DB.Delete(&models.Trigger{}, "id = ?", id).Error
}

// Update implements Repository.
func (r *repo) Update(trigger *schemas.UpdateTriggerRequest, id *uuid.UUID) error {
	// Prepare a map to hold the fields that will be updated
	updateFields := make(map[string]interface{})

	// Add the fields to be updated to the map
	if trigger.ScheduleTime != nil {
		updateFields["schedule_time"] = trigger.ScheduleTime
	}
	if trigger.APIURL != nil {
		updateFields["api_url"] = trigger.APIURL
	}
	if trigger.APIPayload != nil {
		updateFields["api_payload"] = trigger.APIPayload
	}
	if trigger.IntervalSecs != nil {
		updateFields["interval_secs"] = trigger.IntervalSecs
	}
	if trigger.NumberOfOccurrences != nil {
		updateFields["number_of_occurrences"] = trigger.NumberOfOccurrences
	}

	// Update the trigger in the database
	return r.DB.Model(&models.Trigger{}).Where("id = ?", id).Updates(updateFields).Error
}

func (r *repo) FindScheduledTriggers(startTime, endTime time.Time) ([]*models.Trigger, error) {
	var triggers []*models.Trigger
	err := r.DB.Where("schedule_time BETWEEN ? AND ? AND execution_status = ?", startTime, endTime, models.Initialized).Find(&triggers).Error
	return triggers, err
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
		ExecutionStatus:     models.Initialized,
	}
	err := r.DB.Create(t).Error
	if err != nil {
		return nil, err
	}
	return &triggerID, nil

}

func (r *repo) UpdateExecutionStatus(id *uuid.UUID, status models.ExecutionStatus) error {
	return r.DB.Model(&models.Trigger{}).Where("id = ?", id).Update("execution_status", status).Error
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
