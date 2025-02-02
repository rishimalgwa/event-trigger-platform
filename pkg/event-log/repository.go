package eventlog

import (
	"time"

	"github.com/google/uuid"
	"github.com/rishimalgwa/event-trigger-platform/pkg/models"
)

type Repository interface {
	Save(eventLog *models.EventLog) error
	ArchiveLogs(before time.Time) error
	DeleteLogs(before time.Time) error
	MarkDeleteLogs(id uuid.UUID) error
}
