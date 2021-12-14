package apiserver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gopherschool/http-rest-api/internal/app/models"
	"github.com/gopherschool/http-rest-api/internal/app/store/teststore"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServer_AuthenticateUser(t *testing.T) {
	store := teststore.NewStore()
	u := models.TestUser(t)
	store.User().Create(u)

	testCases := []struct {
		name         string
		cookieValue  map[interface{}]interface{}
		expectedCode int
	}{
		{
			name: "authenticated",
			cookieValue: map[interface{}]interface{}{
				"user_id": u.ID,
			},
			expectedCode: http.StatusOK,
		},
		{
			name:         "not authenticated",
			cookieValue:  nil,
			expectedCode: http.StatusUnauthorized,
		},
	}

	secretKey := []byte("secret")
	// sending a simple random key to NewCookieStore
	s := newServer(store, sessions.NewCookieStore(secretKey))
	// we need to generate a string and attach it to request header from cookieValue, send it on server
	// and then try to get some session on the server and check whether user exists or not
	// for that, let's use secure cookie
	sc := securecookie.New(secretKey, nil)
	// add fake handler that implements http Handler interface
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/", nil)
			cookieStr, _ := sc.Encode(sessionName, tc.cookieValue) // pass cookie from test-case data
			// Set header like session_name=encrypted_cookie
			req.Header.Set("Cookie", fmt.Sprintf("%s=%s", sessionName, cookieStr))
			s.authenticateUser(handler).ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}

func TestServerHandleUsersCreate(t *testing.T) {
	s := newServer(teststore.NewStore(), sessions.NewCookieStore([]byte("random_secret")))
	testCases := []struct {
		name         string
		payload      interface{}
		expectedCode int
	}{
		{
			name: "valid data",
			payload: map[string]string{
				"email":    "user@example.org",
				"password": "password",
			},
			expectedCode: http.StatusCreated, // see handleUsersCreate() in server.go
		},
		{
			name:         "invalid payload",
			payload:      "invalid_payload",
			expectedCode: http.StatusBadRequest, // see handleUsersCreate() in server.go
		},
		{
			name: "invalid params",
			payload: map[string]string{
				"email": "invalid",
			},
			expectedCode: http.StatusUnprocessableEntity, // see handleUsersCreate() in server.go
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// TODO: find out about `httptest.NewRecorder`
			rec := httptest.NewRecorder()

			// TODO: find out about next 2 lines
			b := &bytes.Buffer{}
			json.NewEncoder(b).Encode(tc.payload)

			// Perform new request to /users endpoint with payload defined above
			req, _ := http.NewRequest(http.MethodPost, "/users", b)
			// TODO: find out about `ServeHTTP`
			s.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}

func TestServerHandleSessionsCreate(t *testing.T) {
	u := models.TestUser(t)
	store := teststore.NewStore()
	store.User().Create(u)
	s := newServer(store, sessions.NewCookieStore([]byte("random_secret")))
	testCases := []struct {
		name         string
		payload      interface{}
		expectedCode int
	}{
		{
			name: "valid data",
			payload: map[string]string{
				"email":    u.Email,
				"password": u.Password,
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "invalid email",
			payload: map[string]string{
				"email":    "invalid",
				"password": "qwe123QWE",
			},
			expectedCode: http.StatusUnauthorized,
		},
		{
			name: "invalid password",
			payload: map[string]string{
				"email":    "invalid",
				"password": "qwe123QWE!@#",
			},
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "invalid payload",
			payload:      "invalid",
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()

			b := &bytes.Buffer{}
			json.NewEncoder(b).Encode(tc.payload)

			req, _ := http.NewRequest(http.MethodPost, "/sessions", b)
			s.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}
