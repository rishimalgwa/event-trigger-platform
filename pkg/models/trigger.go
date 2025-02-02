package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TriggerType string

const (
	ScheduledTrigger TriggerType = "scheduled"
	APITrigger       TriggerType = "api"
)

type Trigger struct {
	BaseModel
	UserID              uuid.UUID   `gorm:"type:uuid;not null"`
	Type                TriggerType `gorm:"type:varchar(20);not null;check:type IN ('scheduled', 'api')"`
	ScheduleTime        *time.Time  `gorm:"default:null"`
	IntervalSecs        *int        `gorm:"default:null"`
	IsRecurring         bool        `gorm:"default:false"`
	APIURL              *string     `gorm:"type:text;default:null"`
	APIPayload          *string     `gorm:"type:jsonb;default:null"`
	NumberOfOccurrences *int        `gorm:"type:int;default:null"` // For recurring triggers
	ExecutedCount       int         `gorm:"type:int;default:0"`    // Tracks executions
}

func (t *Trigger) BeforeCreate(tx *gorm.DB) error {
	t.ID = uuid.New()
	return nil
}
