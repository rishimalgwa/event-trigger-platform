package main

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/spf13/viper"

	"github.com/rishimalgwa/event-trigger-platform/api/cache"
	"github.com/rishimalgwa/event-trigger-platform/api/db"
	"github.com/rishimalgwa/event-trigger-platform/api/kafka"
	"github.com/rishimalgwa/event-trigger-platform/api/migrations"
	"github.com/rishimalgwa/event-trigger-platform/api/router"
	"github.com/rishimalgwa/event-trigger-platform/api/utils"
)

func healthCheck(c *fiber.Ctx) error {
	return c.SendString("OK")
}

func main() {
	// Set global configuration
	utils.ImportEnv()

	// Init redis
	cache.GetRedis()

	// Init Validators
	utils.InitValidators()

	// Create Fiber
	app := fiber.New(fiber.Config{})

	app.Get("/", healthCheck)
	app.Get("/health", healthCheck)

	app.Use(logger.New(logger.Config{Next: func(c *fiber.Ctx) bool {
		return strings.HasPrefix(c.Path(), "api")
	}}))

	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "",
		AllowHeaders: "*",
	}))

	//Connect and migrate the db
	if viper.GetBool("MIGRATE") {
		migrations.Migrate()
	}

	// Initialize DB
	db.InitServices()
	// Start Kafka consumer in a separate goroutine
	go kafka.StartTriggerConsumer(kafka.Consumer, kafka.Producer, db.TriggerSvc, db.EventLogSvc)

	// Mount Routes
	router.MountRoutes(app)

	// Serve static files (HTML, CSS, JS)
	app.Static("/", "./static")

	// Get Port
	port := utils.GetPort()

	// Start Fiber
	err := app.Listen(fmt.Sprintf(":%s", port))
	if err != nil {
		panic(err)
	}

}
