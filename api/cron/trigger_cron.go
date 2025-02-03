package cron

import (
	"encoding/json"
	"log"
	"time"

	"github.com/IBM/sarama"
	"github.com/rishimalgwa/event-trigger-platform/api/constants"
	"github.com/rishimalgwa/event-trigger-platform/api/db"
	"github.com/rishimalgwa/event-trigger-platform/api/kafka"
)

func StartTriggerCron() {
	ticker := time.NewTicker(1 * time.Minute)

	go func() {
		for range ticker.C {
			// Get current time and next minute time
			now := time.Now()
			nextMinute := now.Add(1 * time.Minute)

			// Fetch triggers from DB scheduled within the next 1 minute
			triggers, err := db.TriggerSvc.FindScheduledTriggers(now, nextMinute)
			if err != nil {
				log.Printf("Error fetching triggers: %v", err)
				continue
			}

			for _, trigger := range triggers {
				// Prepare Kafka message
				triggerEvent := map[string]interface{}{
					"triggerID":       trigger.ID,
					"eventTime":       trigger.ScheduleTime.Format(time.RFC3339),
					"isRecurring":     trigger.IsRecurring,
					"intervalSecs":    trigger.IntervalSecs,
					"occurrencesLeft": trigger.NumberOfOccurrences,
				}

				msgBytes, _ := json.Marshal(triggerEvent)
				msg := &sarama.ProducerMessage{
					Topic: constants.KAFKA_SCHEDULED_TRIGGERS_TOPIC,
					Value: sarama.StringEncoder(msgBytes),
				}

				// Send to Kafka
				_, _, err := kafka.Producer.SendMessage(msg)
				if err != nil {
					log.Printf("Error sending trigger to Kafka: %v", err)
				} else {
					log.Printf("Trigger %s sent to Kafka", trigger.ID)
				}
			}
		}
	}()
}
