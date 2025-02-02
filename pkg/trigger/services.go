package trigger

import (
	"bytes"
	"encoding/json"
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
	ScheduleTrigger(trigger *models.Trigger) error
	ExecuteTrigger(trigger *models.Trigger) error
}

type triggerSvc struct {
	repo      Repository
	kafkaProd sarama.SyncProducer
}

func NewService(r Repository, kafkaProd sarama.SyncProducer) Service {
	return &triggerSvc{repo: r, kafkaProd: kafkaProd}
}

// Find implements Service.
func (t *triggerSvc) Find(id *uuid.UUID) (*models.Trigger, error) {
	return t.repo.Find(id)
}

// Save implements Service.
func (t *triggerSvc) Save(trigger *schemas.CreateTriggerRequest) (*uuid.UUID, error) {
	return t.repo.Save(trigger)
}

func (t *triggerSvc) ScheduleTrigger(trigger *models.Trigger) error {
	// Determine execution time
	var executeAt time.Time
	if trigger.ScheduleTime != nil {
		executeAt = *trigger.ScheduleTime
	} else {
		executeAt = time.Now().Add(time.Duration(*trigger.IntervalSecs) * time.Second)
	}

	// Prepare Kafka message
	triggerEvent := map[string]interface{}{
		"triggerID":       trigger.ID,
		"eventTime":       executeAt.Format(time.RFC3339),
		"isRecurring":     trigger.IsRecurring,
		"intervalSecs":    trigger.IntervalSecs,
		"occurrencesLeft": trigger.NumberOfOccurrences, // Track recurring execution count
	}

	msgBytes, _ := json.Marshal(triggerEvent)
	msg := &sarama.ProducerMessage{
		Topic: "scheduled-triggers",
		Value: sarama.StringEncoder(msgBytes),
	}

	// Send to Kafka
	_, _, err := t.kafkaProd.SendMessage(msg)
	println("Message sent to Kafka")
	return err
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
