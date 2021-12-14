package models_test

import (
	"github.com/gopherschool/http-rest-api/internal/app/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUser_Validate(t *testing.T) {
	testCases := []struct {
		name    string
		u       func() *models.User // function which returns user on output. Needed for user update during tests
		isValid bool
	}{
		{
			name: "valid",
			u: func() *models.User {
				return models.TestUser(t)
			},
			isValid: true,
		},

		{
			name: "empty email",
			u: func() *models.User {
				u := models.TestUser(t)
				u.Email = ""
				return u
			},
			isValid: false,
		},

		{
			name: "invalid email",
			u: func() *models.User {
				u := models.TestUser(t)
				u.Email = "invalid"
				return u
			},
			isValid: false,
		},

		{
			name: "empty password",
			u: func() *models.User {
				u := models.TestUser(t)
				u.Password = ""
				return u
			},
			isValid: false,
		},

		{
			name: "short password",
			u: func() *models.User {
				u := models.TestUser(t)
				u.Password = "123"

				return u
			},
			isValid: false,
		},

		{
			name: "with encrypted password",
			u: func() *models.User {
				u := models.TestUser(t)
				u.Password = ""
				u.EncryptedPassword = "12345"
				return u
			},
			isValid: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.isValid {
				assert.NoError(t, tc.u().Validate())
			} else {
				assert.Error(t, tc.u().Validate())
			}
		})
	}

}

func TestUser_BeforeCreate(t *testing.T) {
	u := models.TestUser(t)
	assert.NoError(t, u.BeforeCreate())
	assert.NotEmpty(t, u.EncryptedPassword)
}
