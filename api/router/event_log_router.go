package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rishimalgwa/event-trigger-platform/api/controllers"
)

func MountEventLogRoutes(apiGroup fiber.Router) {
	thisController := controllers.EventLogController{}
	apiGroup.Get("/eventlog", thisController.GetEventLogs)
}
