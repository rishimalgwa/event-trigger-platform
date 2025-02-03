package models

import (
	"time"

	"github.com/google/uuid"
)

type TriggerType string

const (
	ScheduledTrigger TriggerType = "scheduled"
	APITrigger       TriggerType = "api"
)

type ExecutionStatus string

const (
	Executed    ExecutionStatus = "executed"
	Failed      ExecutionStatus = "failed"
	Executing   ExecutionStatus = "executing"
	Initialized ExecutionStatus = "initialized"
)

type Trigger struct {
	BaseModel
	UserID              uuid.UUID       `gorm:"type:uuid;not null"`
	Type                TriggerType     `gorm:"type:varchar(20);not null;check:type IN ('scheduled', 'api')"`
	ScheduleTime        *time.Time      `gorm:"default:null"`
	IntervalSecs        *int            `gorm:"default:null"`
	IsRecurring         *bool           `gorm:"default:false"`
	APIURL              *string         `gorm:"type:text;default:null"`
	APIPayload          *string         `gorm:"type:jsonb;default:null"`
	NumberOfOccurrences *int            `gorm:"type:int;default:null"` // For recurring triggers
	ExecutedCount       int             `gorm:"type:int;default:0"`    // Tracks executions
	ExecutionStatus     ExecutionStatus `gorm:"type:varchar(20);default:'initialized'; check:execution_status IN ('initialized', 'executing', 'executed', 'failed')"`
}
