package sqlstore

import (
	"database/sql"

	_ "github.com/lib/pq" // Anonymous import to skip import of methods

	"github.com/gopherschool/http-rest-api/internal/app/store"
)

type Store struct {
	db             *sql.DB
	userRepository *UserRepository
}

// NewStore returns pointer on store
func NewStore(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

// User is special method to avoid using repositories without the store.
// Example of such call: store.User().Create()
func (s *Store) User() store.UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = &UserRepository{store: s}
	return s.userRepository
}
