package store

import "github.com/gopherschool/http-rest-api/internal/app/models"

// UserRepository is an interface for user repositories
type UserRepository interface {
	Create(*models.User) error
	FindByEmail(string) (*models.User, error)
	FindByID(int) (*models.User, error)
}
