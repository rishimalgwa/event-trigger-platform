package users

import (
	"github.com/google/uuid"
	"github.com/rishimalgwa/go-template/pkg/models"
)

type Service interface{}

type userSvc struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &userSvc{repo: r}
}

func (u *userSvc) Find(id *uuid.UUID) (*models.Users, error) {
	return u.repo.Find(id)
}
