package kafka

import (
	"encoding/json"
	"log"
	"time"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
	"github.com/rishimalgwa/event-trigger-platform/api/constants"
	eventlog "github.com/rishimalgwa/event-trigger-platform/pkg/event-log"
	"github.com/rishimalgwa/event-trigger-platform/pkg/models"
	"github.com/rishimalgwa/event-trigger-platform/pkg/trigger"
)

// StartTriggerConsumer listens to the Kafka topic and executes triggers
func StartTriggerConsumer(kafkaConsumer sarama.Consumer, kafkaProducer sarama.SyncProducer, triggerSvc trigger.Service, eventLogSvc eventlog.Service) {
	part, err := kafkaConsumer.ConsumePartition(constants.KAFKA_SCHEDULED_TRIGGERS_TOPIC, 0, sarama.OffsetNewest)
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
		isTest, _ := event["isTest"].(bool) // Check for the test flag

		execTime, _ := time.Parse(time.RFC3339, executeTime)

		if isTest {
			apiPayload, _ := event["apiPayload"].(string)
			apiURL, _ := event["apiURL"].(string)
			log.Printf("Test Trigger received, executing once")
			testTriggerID := uuid.New()
			// Sleep until execution time for scheduled triggers
			time.Sleep(time.Until(execTime))
			// Log before execution
			eventLog := &models.EventLog{
				TriggerID:   testTriggerID,
				TriggeredAt: time.Now(),
				Status:      models.Active,
				APIPayload:  &apiPayload,
				APIURL:      &apiURL,
				IsManual:    true,
			}

			// Save event log for the test trigger
			eventLogSvc.SaveEventLog(eventLog)

			// Execute the test trigger

			err = triggerSvc.ExecuteTrigger(&models.Trigger{
				BaseModel: models.BaseModel{
					ID: testTriggerID,
				},
				Type:       models.APITrigger,
				APIURL:     &apiURL,
				APIPayload: &apiPayload,
			})
			if err != nil {
				triggerSvc.UpdateExecutionStatus(&testTriggerID, models.Failed)
				log.Printf("Trigger execution failed: %v", err)
				continue
			}

			// Mark as Executed for test triggers
			triggerSvc.UpdateExecutionStatus(&testTriggerID, models.Executed)
			println("Test Trigger executed")

			// Skip the recurring logic as test triggers are not recurring
			continue
		}

		// Sleep until execution time for regular scheduled triggers
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
				Topic: constants.KAFKA_SCHEDULED_TRIGGERS_TOPIC,
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
