package kafka

import (
	"encoding/json"
	"log"
	"time"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
	eventlog "github.com/rishimalgwa/event-trigger-platform/pkg/event-log"
	"github.com/rishimalgwa/event-trigger-platform/pkg/models"
	"github.com/rishimalgwa/event-trigger-platform/pkg/trigger"
)

// StartTriggerConsumer listens to the Kafka topic and executes triggers
func StartTriggerConsumer(kafkaConsumer sarama.Consumer, kafkaProducer sarama.SyncProducer, triggerSvc trigger.Service, eventLogSvc eventlog.Service) {
	part, err := kafkaConsumer.ConsumePartition("scheduled-triggers", 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("Error starting Kafka consumer: %v", err)
	}

	log.Println("Kafka Trigger Consumer started...")

	for msg := range part.Messages() {
		var event map[string]interface{}
		err := json.Unmarshal(msg.Value, &event)
		if err != nil {
			log.Printf("Error decoding event: %v", err)
			continue
		}

		triggerID, _ := event["triggerID"].(string)
		executeTime, _ := event["eventTime"].(string)
		isRecurring, _ := event["isRecurring"].(bool)
		intervalSecs, _ := event["intervalSecs"].(float64)
		occurrencesLeft, _ := event["occurrencesLeft"].(float64)

		// Convert execution time to Go time format
		execTime, _ := time.Parse(time.RFC3339, executeTime)

		// Sleep until execution time
		time.Sleep(time.Until(execTime))

		// Find the trigger from DB
		triggerUUID, _ := uuid.Parse(triggerID)
		trigger, err := triggerSvc.Find(&triggerUUID)
		if err != nil {
			log.Printf("Trigger not found: %v", err)
			continue
		}

		// Log before execution
		eventLog := &models.EventLog{
			TriggerID:   trigger.ID,
			TriggeredAt: time.Now(),
			Status:      models.Active,
			APIPayload:  trigger.APIPayload,
			APIURL:      trigger.APIURL,
		}
		triggerSvc.UpdateExecutionStatus(&trigger.ID, models.Executing)

		// Execute trigger
		err = triggerSvc.ExecuteTrigger(trigger)
		if err != nil {
			triggerSvc.UpdateExecutionStatus(&trigger.ID, models.Failed)
			log.Printf("Trigger execution failed: %v", err)
			continue
		}

		// Save event log
		eventLogSvc.SaveEventLog(eventLog)

		// Handle recurring triggers properly
		if isRecurring && occurrencesLeft > 1 {
			event["eventTime"] = time.Now().Add(time.Duration(intervalSecs) * time.Second).Format(time.RFC3339)
			event["occurrencesLeft"] = occurrencesLeft - 1

			msgBytes, _ := json.Marshal(event)
			msg := &sarama.ProducerMessage{
				Topic: "scheduled-triggers",
				Value: sarama.StringEncoder(msgBytes),
			}
			_, _, _ = kafkaProducer.SendMessage(msg)
		} else {
			// Only mark as Executed when it's NOT recurring or all occurrences are done
			triggerSvc.UpdateExecutionStatus(&trigger.ID, models.Executed)
		}

		println("Trigger executed")
	}
}
