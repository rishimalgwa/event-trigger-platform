package trigger

import (
	"github.com/google/uuid"
	"github.com/rishimalgwa/event-trigger-platform/api/schemas"
	"github.com/rishimalgwa/event-trigger-platform/pkg/models"
)

type Repository interface {
	Find(id *uuid.UUID) (*models.Trigger, error)
	Save(trigger *schemas.CreateTriggerRequest) error
}
