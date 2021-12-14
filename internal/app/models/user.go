package models

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"golang.org/x/crypto/bcrypt"
)

// User models doesn't know anything about interaction with DB
// Repositories will be responsible for this kind of interaction
type User struct {
	ID                int    `json:"id"`
	Email             string `json:"email"`
	Password          string `json:"password,omitempty"` // is password is empty, then don't return it
	EncryptedPassword string `json:"-"`                  // do not render encr password
}

func (u *User) Validate() error {
	// check that Email is required and it has correct form
	// custom validator was implemented for password check
	return validation.ValidateStruct(
		u,
		validation.Field(&u.Email, validation.Required, is.Email),
		validation.Field(&u.Password, validation.By(requiredIf(u.EncryptedPassword == "")), validation.Length(6, 100)),
	)
}

// Lesson3, timeframe 1:10
func (u *User) BeforeCreate() error {
	if len(u.Password) > 0 {
		enc, err := encryptString(u.Password)
		if err != nil {
			return err
		}

		u.EncryptedPassword = enc
	}
	return nil
}

// Sanitize redefines private attributes that shouldn't be available outside
func (u *User) Sanitize() {
	u.Password = ""
}

// ComparePasswords check that password from session request corresponds to encrypted
// will return true if comparison is OK
func (u *User) ComparePasswords(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.EncryptedPassword), []byte(password)) == nil
}

func encryptString(s string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.MinCost)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
