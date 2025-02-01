package users

import (
	"github.com/google/uuid"
	"github.com/rishimalgwa/event-trigger-platform/pkg/models"
)

type Repository interface {
	Find(id *uuid.UUID) (*models.Users, error)
}
