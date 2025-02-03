package trigger

import (
	"bytes"
	"errors"
	"net/http"
	"time"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
	"github.com/rishimalgwa/event-trigger-platform/api/schemas"
	"github.com/rishimalgwa/event-trigger-platform/pkg/models"
)

type Service interface {
	Find(id *uuid.UUID) (*models.Trigger, error)
	Save(trigger *schemas.CreateTriggerRequest) (*uuid.UUID, error)
	UpdateExecutionStatus(id *uuid.UUID, status models.ExecutionStatus) error
	Update(trigger *schemas.UpdateTriggerRequest, id *uuid.UUID) error
	ExecuteTrigger(trigger *models.Trigger) error
	UpdateTrigger(triggerID uuid.UUID, updates *schemas.UpdateTriggerRequest) error
	DeleteTrigger(triggerID uuid.UUID) error
	FindScheduledTriggers(start, end time.Time) ([]*models.Trigger, error)
	FindAll() ([]*models.Trigger, error)
}

type triggerSvc struct {
	repo      Repository
	kafkaProd sarama.SyncProducer
}

func NewService(r Repository, kafkaProd sarama.SyncProducer) Service {
	return &triggerSvc{repo: r, kafkaProd: kafkaProd}
}

func (t *triggerSvc) FindAll() ([]*models.Trigger, error) {
	return t.repo.FindAll()
}

// Update implements Service.
func (t *triggerSvc) Update(trigger *schemas.UpdateTriggerRequest, id *uuid.UUID) error {
	return t.repo.Update(trigger, id)
}

func (t *triggerSvc) FindScheduledTriggers(start, end time.Time) ([]*models.Trigger, error) {
	return t.repo.FindScheduledTriggers(start, end)
}

// Find implements Service.
func (t *triggerSvc) Find(id *uuid.UUID) (*models.Trigger, error) {
	return t.repo.Find(id)
}

// Save implements Service.
func (t *triggerSvc) Save(trigger *schemas.CreateTriggerRequest) (*uuid.UUID, error) {
	return t.repo.Save(trigger)
}

func (t *triggerSvc) UpdateExecutionStatus(id *uuid.UUID, status models.ExecutionStatus) error {
	return t.repo.UpdateExecutionStatus(id, status)
}

func (t *triggerSvc) ExecuteTrigger(trigger *models.Trigger) error {

	// Handle API trigger execution
	if trigger.APIURL != nil {

		resp, err := http.Post(*trigger.APIURL, "application/json", bytes.NewBuffer([]byte(*trigger.APIPayload)))
		if err != nil {

			return err
		}
		println(resp.StatusCode)
		defer resp.Body.Close()
	}

	return nil
}
func (t *triggerSvc) UpdateTrigger(triggerID uuid.UUID, updates *schemas.UpdateTriggerRequest) error {
	trigger, err := t.repo.Find(&triggerID)
	if err != nil {
		return errors.New("trigger not found")
	}

	if trigger.ExecutionStatus != models.Initialized {
		return errors.New("cannot update a trigger that is already executing or executed")
	}

	// Save updated trigger
	return t.repo.Update(updates, &triggerID)
}

func (t *triggerSvc) DeleteTrigger(triggerID uuid.UUID) error {
	trigger, err := t.repo.Find(&triggerID)
	if err != nil {
		return errors.New("trigger not found")
	}

	if trigger.ExecutionStatus != models.Initialized {
		return errors.New("cannot delete a trigger that is already executing or executed")
	}

	// Delete trigger from DB
	return t.repo.Delete(&triggerID)
}
