package teststore

import (
	"github.com/gopherschool/http-rest-api/internal/app/models"
	"github.com/gopherschool/http-rest-api/internal/app/store"
)

// another realization of store for tests (?)

type Store struct {
	userRepository *UserRepository
}

// NewStore returns pointer on store
func NewStore() *Store {
	return &Store{}
}

// User is special method to avoid using repositories without the store.
// Example of such call: store.User().Create()
func (s *Store) User() store.UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = &UserRepository{
		store: s,
		users: make(map[int]*models.User),
	}

	return s.userRepository
}
