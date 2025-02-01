package db

import "github.com/rishimalgwa/go-template/pkg/users"

var (
	UsersSvc users.Service = nil
)

func InitServices() {
	db := GetDB()

	usersRepo := users.NewPostgresRepo(db)
	UsersSvc = users.NewService(usersRepo)
}
