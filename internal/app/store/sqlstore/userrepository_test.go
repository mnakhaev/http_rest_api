package sqlstore_test

import (
	"github.com/gopherschool/http-rest-api/internal/app/models"
	"github.com/gopherschool/http-rest-api/internal/app/store"
	"github.com/gopherschool/http-rest-api/internal/app/store/sqlstore"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

// tests for Create method
func TestUserRepository_Create(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("users") // cleaning users table

	s := sqlstore.NewStore(db)
	u := models.TestUser(t)
	u.ID = rand.Intn(1000)
	assert.NoError(t, s.User().Create(u)) // check that no error raised
	assert.NotNil(t, u)                   // check that user is not nil
}

func TestUserRepository_FindByEmail(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("users") // cleaning users table

	s := sqlstore.NewStore(db)
	email := "user123@example.org"
	_, err := s.User().FindByEmail(email)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	u := models.TestUser(t)
	u.Email = email
	s.User().Create(u)
	u, err = s.User().FindByEmail(email)
	assert.NoError(t, err)
	assert.NotNil(t, u)
}

func TestUserRepository_FindByID(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("users")

	s := sqlstore.NewStore(db)
	u1 := models.TestUser(t)
	if err := s.User().Create(u1); err != nil {
		t.Fatal(err)
	}
	u2 := models.TestUser(t)

	u2, err := s.User().FindByID(u1.ID)
	assert.NoError(t, err)
	assert.NotNil(t, u2)
	assert.Equal(t, u2.ID, u1.ID)
}
