package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EventStatus string

const (
	Active   EventStatus = "active"   // First 2 hours
	Archived EventStatus = "archived" // Next 46 hours
	Deleted  EventStatus = "deleted"  // After 48 hours
)

type EventLog struct {
	BaseModel
	TriggerID   uuid.UUID   `gorm:"type:uuid;not null"`
	TriggeredAt time.Time   `gorm:"autoCreateTime"`
	Status      EventStatus `gorm:"type:varchar(20);not null;check:status IN ('active', 'archived', 'deleted')"`
	APIPayload  *string     `gorm:"type:jsonb;default:null"`
	IsManual    bool        `gorm:"default:false"`
}

func (e *EventLog) BeforeCreate(tx *gorm.DB) error {
	e.ID = uuid.New()
	return nil
}
