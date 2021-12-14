package teststore

import (
	"github.com/gopherschool/http-rest-api/internal/app/models"
	"github.com/gopherschool/http-rest-api/internal/app/store"
)

// UserRepository structure for tests
type UserRepository struct {
	store *Store
	users map[int]*models.User
}

// Create test user in `users` map
func (r *UserRepository) Create(u *models.User) error {
	// check if user is valid. if OK - run BeforeCreate callback
	if err := u.Validate(); err != nil {
		return err
	}

	if err := u.BeforeCreate(); err != nil {
		return err
	}

	u.ID = len(r.users) + 1
	r.users[u.ID] = u

	return nil
}

// FindByEmail in `users` map
func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	for _, u := range r.users {
		if u.Email == email {
			return u, nil
		}
	}
	return nil, store.ErrRecordNotFound
}

// TODO: implement till the end
// FindByID in `users` map
func (r *UserRepository) FindByID(ID int) (*models.User, error) {
	u, ok := r.users[ID]
	if !ok {
		return nil, store.ErrRecordNotFound
	}

	return u, nil
}
