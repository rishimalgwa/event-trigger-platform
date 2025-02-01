package migrations

import (
	"github.com/rishimalgwa/event-trigger-platform/api/db"
	"github.com/rishimalgwa/event-trigger-platform/pkg/models"
)

func Migrate() {
	database := db.GetDB()
	database.Raw("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
	database.AutoMigrate(&models.Users{}, &models.Items{})
}
