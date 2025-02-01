package db

import (
	eventlog "github.com/rishimalgwa/event-trigger-platform/pkg/event-log"
	"github.com/rishimalgwa/event-trigger-platform/pkg/trigger"
)

var (
	TriggerSvc  trigger.Service  = nil
	EventLogSvc eventlog.Service = nil
)

func InitServices() {
	db := GetDB()

	triggerRepo := trigger.NewPostgresRepo(db)
	TriggerSvc = trigger.NewService(triggerRepo)

	eventLogRepo := eventlog.NewPostgresRepo(db)
	EventLogSvc = eventlog.NewService(eventLogRepo)
}
