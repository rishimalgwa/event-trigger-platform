package users

import (
	"github.com/google/uuid"
	"github.com/rishimalgwa/go-template/pkg/models"
)

type Repository interface {
	Find(id *uuid.UUID) (*models.Users, error)
}
