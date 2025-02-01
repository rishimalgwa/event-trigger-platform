package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rishimalgwa/event-trigger-platform/api/views"
)

type EventLogController struct{}

func (e *EventLogController) GetEventLogs(ctx *fiber.Ctx) error {
	return views.OK(ctx, nil)
}
