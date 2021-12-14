package teststore_test

import (
	"github.com/gopherschool/http-rest-api/internal/app/models"
	"github.com/gopherschool/http-rest-api/internal/app/store"
	"github.com/gopherschool/http-rest-api/internal/app/store/teststore"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

// tests for Create method
func TestUserRepository_Create(t *testing.T) {
	s := teststore.NewStore()
	u := models.TestUser(t)
	u.ID = rand.Intn(1000)
	assert.NoError(t, s.User().Create(u)) // check that no error raised
	assert.NotNil(t, u)                   // check that user is not nil
}

func TestUserRepository_FindByEmail(t *testing.T) {
	s := teststore.NewStore()
	u1 := models.TestUser(t)
	_, err := s.User().FindByEmail(u1.Email)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	s.User().Create(u1)
	u2, err := s.User().FindByEmail(u1.Email)
	assert.NoError(t, err)
	assert.NotNil(t, u2)
}

func TestUserRepository_FindByID(t *testing.T) {
	s := teststore.NewStore()
	u1 := models.TestUser(t)
	_, err := s.User().FindByID(u1.ID)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	s.User().Create(u1)
	u2, err := s.User().FindByID(u1.ID)
	assert.NoError(t, err)
	assert.NotNil(t, u2)
}
