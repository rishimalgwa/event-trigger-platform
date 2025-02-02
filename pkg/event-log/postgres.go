package eventlog

import (
	"time"

	"github.com/google/uuid"
	"github.com/rishimalgwa/event-trigger-platform/pkg/models"
	"gorm.io/gorm"
)

type repo struct {
	DB *gorm.DB
}

// MarkDeleteLogs implements Repository.
func (r *repo) MarkDeleteLogs(id uuid.UUID) error {
	return r.DB.Model(&models.EventLog{}).Where("id = ?", id).Update("status", "deleted").Error
}

func NewPostgresRepo(db *gorm.DB) Repository {
	return &repo{
		DB: db,
	}
}

func (r *repo) Save(eventLog *models.EventLog) error {
	return r.DB.Create(eventLog).Error
}

func (r *repo) ArchiveLogs(before time.Time) error {
	return r.DB.Model(&models.EventLog{}).Where("created_at < ?", before).Update("status", "archived").Error
}

func (r *repo) DeleteLogs(before time.Time) error {
	return r.DB.Model(&models.EventLog{}).Where("created_at < ?", before).Update("status", "deleted").Error
}
