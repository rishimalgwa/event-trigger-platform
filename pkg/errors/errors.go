package errors

import (
	err "errors"

	"github.com/lib/pq"
)

var (
	ErrSomethingWentWrong = err.New("something went wrong")
	ErrSuperUser          = err.New("cannot remove a superuser")
	ErrIsNotSuperUser     = err.New("user is not superuser")
	ErrIsAlreadyExists    = err.New("already exists")
	ErrRegisteredAgentErr = err.New("registered agent not found")
	ErrUserUnauthorized   = err.New("user is unauthorized")
)

func IsAlreadyExists(thisErr error) bool {
	pqErr, ok := thisErr.(*pq.Error)
	if !ok {
		return false
	}
	if pqErr.Code.Name() == "unique_violation" {
		return true
	}
	return false
}
