package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rishimalgwa/event-trigger-platform/api/controllers"
)

func MountTriggerRoutes(apiGroup fiber.Router) {
	thisController := controllers.TriggerController{}
	apiGroup.Post("/trigger", thisController.CreateTrigger)
	apiGroup.Get("/trigger", thisController.GetAllTriggers)
	apiGroup.Get("/trigger/:id", thisController.GetTriggerById)
	apiGroup.Delete("/trigger/:id", thisController.DeleteTrigger)
	apiGroup.Put("/trigger/:id", thisController.UpdateTrigger)
}
