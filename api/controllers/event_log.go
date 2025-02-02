package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rishimalgwa/event-trigger-platform/api/db"
	"github.com/rishimalgwa/event-trigger-platform/api/views"
)

type EventLogController struct{}

func (e *EventLogController) GetActiveEventLogs(ctx *fiber.Ctx) error {
	logs, err := db.EventLogSvc.GetActiveLogs()
	if err != nil {
		return views.InternalServerError(ctx, err)
	}
	return views.OK(ctx, logs)
}
func (e *EventLogController) GetArchivedEventLogs(ctx *fiber.Ctx) error {
	logs, err := db.EventLogSvc.GetArchivedLogs()
	if err != nil {
		return views.InternalServerError(ctx, err)
	}
	return views.OK(ctx, logs)
}
