package migrations

import (
	"github.com/rishimalgwa/go-template/api/db"
	"github.com/rishimalgwa/go-template/pkg/models"
)

func Migrate() {
	database := db.GetDB()
	database.Raw("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
	database.AutoMigrate(&models.Users{}, &models.Items{})
}
