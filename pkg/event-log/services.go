package eventlog

import (
	"encoding/json"
	"log"
	"time"

	"github.com/IBM/sarama"
	"github.com/rishimalgwa/event-trigger-platform/api/cache"
	"github.com/rishimalgwa/event-trigger-platform/pkg/models"
)

type Service interface {
	SaveEventLog(eventLog *models.EventLog) error
	ArchiveAndDeleteLogs() error
	GetActiveLogs() ([]models.EventLog, error)
	GetArchivedLogs() ([]models.EventLog, error)
}

type eventLogSvc struct {
	repo      Repository
	kafkaProd sarama.SyncProducer
}

func NewService(r Repository, kafkaProd sarama.SyncProducer) Service {
	return &eventLogSvc{repo: r, kafkaProd: kafkaProd}
}

// SaveEventLog saves the event log and triggers archive and delete
func (s *eventLogSvc) SaveEventLog(eventLog *models.EventLog) error {
	err := s.repo.Save(eventLog)
	if err != nil {
		return err
	}

	// Convert event log to JSON
	data, _ := json.Marshal(eventLog)

	// get redis
	redis, ctx := cache.GetRedis()

	// Store log in "active" Redis list with 48-hour expiration
	redis.LPush(ctx, "event_logs_active", data)
	redis.LTrim(ctx, "event_logs_active", 0, 99) // Keep the latest 100 logs
	redis.Expire(ctx, "event_logs_active", 48*time.Hour)

	// Start background task for archiving & deletion
	go s.ArchiveAndDeleteLogs()
	return nil
}

// ArchiveAndDeleteLogs archives logs after 2 hours and deletes them after 48 hours
func (s *eventLogSvc) ArchiveAndDeleteLogs() error {
	twoHoursAgo := time.Now().Add(-2 * time.Hour).UTC()
	fortyEightHoursAgo := time.Now().Add(-48 * time.Hour).UTC()

	// get redis
	redis, ctx := cache.GetRedis()

	// Fetch all active logs from Redis
	logs, err := redis.LRange(ctx, "event_logs_active", 0, -1).Result()
	if err != nil {
		log.Printf("Error fetching active logs: %v", err)
		return err
	}

	for _, logStr := range logs {
		var eventLog models.EventLog
		json.Unmarshal([]byte(logStr), &eventLog)

		// Convert event log timestamp to time.Time
		eventTime := eventLog.TriggeredAt

		// Move logs older than 2 hours to archived list
		if eventTime.Before(twoHoursAgo) {
			redis.LPush(ctx, "event_logs_archived", logStr)
			redis.LTrim(ctx, "event_logs_archived", 0, 199) // Keep the latest 200 archived logs
			redis.Expire(ctx, "event_logs_archived", 48*time.Hour)

			// Remove from active logs
			redis.LRem(ctx, "event_logs_active", 1, logStr)
		}

		// Delete logs older than 48 hours from Redis and DB
		if eventTime.Before(fortyEightHoursAgo) {
			redis.LRem(ctx, "event_logs_archived", 1, logStr)
			s.repo.MarkDeleteLogs(eventLog.ID)
		}
	}

	return nil
}
func (s *eventLogSvc) GetActiveLogs() ([]models.EventLog, error) {
	// get redis
	redis, ctx := cache.GetRedis()
	logs, err := redis.LRange(ctx, "event_logs_active", 0, 99).Result()
	if err != nil {
		return nil, err
	}

	var eventLogs []models.EventLog
	for _, logStr := range logs {
		var eventLog models.EventLog
		json.Unmarshal([]byte(logStr), &eventLog)
		eventLogs = append(eventLogs, eventLog)
	}
	return eventLogs, nil
}

func (s *eventLogSvc) GetArchivedLogs() ([]models.EventLog, error) {
	// get redis
	redis, ctx := cache.GetRedis()
	logs, err := redis.LRange(ctx, "event_logs_archived", 0, 199).Result()
	if err != nil {
		return nil, err
	}

	var eventLogs []models.EventLog
	for _, logStr := range logs {
		var eventLog models.EventLog
		json.Unmarshal([]byte(logStr), &eventLog)
		eventLogs = append(eventLogs, eventLog)
	}
	return eventLogs, nil
}
