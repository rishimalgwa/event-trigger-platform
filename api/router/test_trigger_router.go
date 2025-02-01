package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rishimalgwa/event-trigger-platform/api/controllers"
)

func MountTriggerTestRoutes(apiGroup fiber.Router) {
	thisController := controllers.TestTriggerController{}
	apiGroup.Post("/trigger/test", thisController.TestTrigger)
}
