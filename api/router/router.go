package router

import "github.com/gofiber/fiber/v2"

func MountRoutes(c *fiber.App) {
	apiGroup := c.Group("/api")
	MountTriggerRoutes(apiGroup)
	MountEventLogRoutes(apiGroup)
	MountTriggerTestRoutes(apiGroup)
}
