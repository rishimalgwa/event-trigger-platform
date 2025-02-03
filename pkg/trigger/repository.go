package trigger

import (
	"time"

	"github.com/google/uuid"
	"github.com/rishimalgwa/event-trigger-platform/api/schemas"
	"github.com/rishimalgwa/event-trigger-platform/pkg/models"
)

type Repository interface {
	Find(id *uuid.UUID) (*models.Trigger, error)
	Save(trigger *schemas.CreateTriggerRequest) (*uuid.UUID, error)
	UpdateExecutionStatus(id *uuid.UUID, status models.ExecutionStatus) error
	Update(trigger *schemas.UpdateTriggerRequest, id *uuid.UUID) error
	Delete(id *uuid.UUID) error
	FindScheduledTriggers(start, end time.Time) ([]*models.Trigger, error)
	FindAll() ([]*models.Trigger, error)
}
