package db

import (
	"log"

	"github.com/rishimalgwa/event-trigger-platform/api/kafka"
	eventlog "github.com/rishimalgwa/event-trigger-platform/pkg/event-log"
	"github.com/rishimalgwa/event-trigger-platform/pkg/trigger"
)

var (
	TriggerSvc  trigger.Service  = nil
	EventLogSvc eventlog.Service = nil
)

func InitServices() {
	// Initialize Kafka
	err := kafka.InitializeKafka()
	if err != nil {
		log.Fatal("Failed to initialize Kafka:", err)
	}

	// Initialize Database
	db := GetDB()

	triggerRepo := trigger.NewPostgresRepo(db)
	TriggerSvc = trigger.NewService(triggerRepo, kafka.Producer)

	eventLogRepo := eventlog.NewPostgresRepo(db)
	EventLogSvc = eventlog.NewService(eventLogRepo, kafka.Producer)
}
