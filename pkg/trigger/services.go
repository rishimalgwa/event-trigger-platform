package trigger

import (
	"github.com/google/uuid"
	"github.com/rishimalgwa/event-trigger-platform/pkg/models"
)

type Service interface{}

type triggerSvc struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &triggerSvc{repo: r}
}

func (t *triggerSvc) Find(id *uuid.UUID) (*models.Trigger, error) {
	return t.repo.Find(id)
}
