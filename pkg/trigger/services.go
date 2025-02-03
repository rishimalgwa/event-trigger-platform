package trigger

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
	"github.com/rishimalgwa/event-trigger-platform/api/constants"
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
	ProduceTestTrigger(trigger schemas.CreateTriggerRequest) error
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

func (s *triggerSvc) ProduceTestTrigger(trigger schemas.CreateTriggerRequest) error {
	eventTime, err := time.Parse(time.RFC3339, *trigger.ScheduleTime)
	if err != nil {
		return fmt.Errorf("failed to parse schedule time: %v", err)
	}
	triggerEvent := map[string]interface{}{

		"eventTime":       eventTime,
		"isRecurring":     trigger.IsRecurring,
		"intervalSecs":    trigger.IntervalSecs,
		"occurrencesLeft": trigger.NumberOfOccurrences,
		"apiURL":          trigger.APIURL,
		"apiPayload":      trigger.APIPayload,
		"isTest":          true, // Mark it as a test trigger
	}

	msgBytes, err := json.Marshal(triggerEvent)
	if err != nil {
		return fmt.Errorf("failed to marshal trigger: %v", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: constants.KAFKA_SCHEDULED_TRIGGERS_TOPIC,
		Value: sarama.StringEncoder(msgBytes),
	}

	// Send to Kafka
	_, _, err = s.kafkaProd.SendMessage(msg)
	if err != nil {
		log.Printf("Error sending test trigger to Kafka: %v", err)
		return err
	}

	log.Printf("Test trigger %s sent to Kafka", trigger.Type)
	return nil
}
