package sqlstore

import (
	"database/sql"

	"github.com/gopherschool/http-rest-api/internal/app/models"
	"github.com/gopherschool/http-rest-api/internal/app/store"
)

type UserRepository struct {
	store *Store
}

// Create accepts and returns needed model
func (r *UserRepository) Create(u *models.User) error {
	// check if user is valid. if OK - run BeforeCreate callback
	if err := u.Validate(); err != nil {
		return err
	}

	if err := u.BeforeCreate(); err != nil {
		return err
	}

	// postgres doesn't return IDs by default, but we need to get this ID for successfully created user
	// this ID will be used later somehow
	// Scan method is used to map returned string to passed arguments (should be pointers!)
	return r.store.db.QueryRow(
		"INSERT INTO users (email, encrypted_password) VALUES ($1, $2) RETURNING id",
		u.Email,
		u.EncryptedPassword,
	).Scan(&u.ID)
}

// FindByEmail method is needed for authorization to find user
func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	u := &models.User{}
	// QueryRow returns only one result
	// Scan fills user with data (?) - need to check the docs
	if err := r.store.db.QueryRow(
		"SELECT id, email, encrypted_password FROM users WHERE email = $1",
		email,
	).Scan(
		&u.ID,
		&u.Email,
		&u.EncryptedPassword,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}
	return u, nil
}

func (r *UserRepository) FindByID(id int) (*models.User, error) {
	u := &models.User{}
	if err := r.store.db.QueryRow(
		"SELECT id, email, encrypted_password FROM users WHERE id = $1",
		id,
	).Scan(
		&u.ID,
		&u.Email,
		&u.EncryptedPassword,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}
	return u, nil
}
