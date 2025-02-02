package trigger

import (
	"github.com/google/uuid"
	"github.com/rishimalgwa/event-trigger-platform/api/schemas"
	"github.com/rishimalgwa/event-trigger-platform/pkg/models"
)

type Service interface {
	Find(id *uuid.UUID) (*models.Trigger, error)
	Save(trigger *schemas.CreateTriggerRequest) error
}

type triggerSvc struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &triggerSvc{repo: r}
}

func (t *triggerSvc) Find(id *uuid.UUID) (*models.Trigger, error) {
	return t.repo.Find(id)
}

func (t *triggerSvc) Save(trigger *schemas.CreateTriggerRequest) error {
	return t.repo.Save(trigger)
}
