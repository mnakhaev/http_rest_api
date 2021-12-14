package apiserver

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"

	"github.com/gopherschool/http-rest-api/internal/app/models"
	"github.com/gopherschool/http-rest-api/internal/app/store"
)

// handle incoming requests and realize HTTP handler interface

const (
	sessionName        = "simple_session_name" // will be returned as response cookie - Set-Cookie: simple_session_name=MTN..
	ctxKeyUser  ctxKey = iota                  // TODO: what is iota?
	ctxKeyRequestID
)

// create common error for both wrong email and password (more secire way)
var (
	errIncorrectEmailOrPassword = errors.New("incorrect email or password")
	errNotAuthenticated         = errors.New("not authenticated")
)

type ctxKey int8

type server struct {
	router       *mux.Router
	logger       *logrus.Logger
	store        store.Store    // it's an interface
	sessionStore sessions.Store // gorilla session. Will be returned as response cookie
}

// newServer accepts store interface
func newServer(store store.Store, sessionStore sessions.Store) *server {
	s := &server{
		router:       mux.NewRouter(),
		logger:       logrus.New(),
		store:        store,
		sessionStore: sessionStore,
	}
	s.configureRouter()
	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) configureRouter() {
	// added middleware that sets request ID at the beginning
	s.router.Use(s.setRequestID)
	s.router.Use(s.logRequest)
	// Allow requests from all sources. Response will contain headers "Access-Control-Allow-Origin: *"
	s.router.Use(handlers.CORS(handlers.AllowedOrigins([]string{"*"})))
	s.router.HandleFunc("/users", s.handleUsersCreate()).Methods("POST")
	// Create new session for user. Will be returned as response header
	s.router.HandleFunc("/sessions", s.handleSessionsCreate()).Methods("POST")

	// add new sub-router that will be hidden by middleware and will ask user for authentication
	// middleware will work with URLs like /private/***
	// add user to context, gets user in handler and renders it
	private := s.router.PathPrefix("/private").Subrouter()
	private.Use(s.authenticateUser)
	private.HandleFunc("/whoami", s.handleWhoami()).Methods("GET")
}

// setRequestID middleware will set unique ID for every input request that will be returned in header and used inside of our system
func (s *server) setRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New().String()
		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyRequestID, id)))
	})
}

// logRequest ...
func (s *server) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := s.logger.WithFields(logrus.Fields{
			"remote_addr": r.RemoteAddr,
			"request_id":  r.Context().Value(ctxKeyRequestID),
		})
		// started GET /endpoint
		logger.Infof("started %s %s", r.Method, r.RequestURI)

		start := time.Now() // catch time when request handling was started
		customWriter := &responseWriter{w, http.StatusOK}
		next.ServeHTTP(customWriter, r)
		// cannot get status code from current response writer
		// we will define our own response writer to workaround the problem

		logger.Infof(
			"completed with %d %s in %v",
			customWriter.code,
			http.StatusText(customWriter.code),
			time.Now().Sub(start),
		)
	})
}

// authenticateUser accept next handler/middleware
func (s *server) authenticateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// firstly, get current user's session from its request
		session, err := s.sessionStore.Get(r, sessionName)
		if err != nil {
			// return 500 because it's our error
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		// get user ID from session. `Values` is one of session's parameters
		id, ok := session.Values["user_id"]
		if !ok {
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		u, err := s.store.User().FindByID(id.(int))
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		// if user was found, then the request is considered as authenticated
		// then, next handler is called

		// we don't want to repeat this operations again if we need to make any request for already authenticated user.
		// to avoid this, we will attach authenticated user to the context of current request
		// context.WithValue accepts 3 args: parent Context, key, val interface{}
		// r.Context() is a context for current request and this is a parent context
		// key is a context key. It's recommended to create a new type for context keys
		// value is a user
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyUser, u)))
	})
}

// handleWhoami renders user that will be taken from context
// we assume here that user is already logged in and we have written him into context and can make a call to him
func (s *server) handleWhoami() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// transform context key user to *models.User type
		s.respond(w, r, http.StatusOK, r.Context().Value(ctxKeyUser).(*models.User))
	}
}

// Check if this function returns pointer
func (s *server) handleUsersCreate() http.HandlerFunc {
	// request describes parameters needed for authentication of user
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			// error accepts response writer, request and response code (400) and error
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		u := &models.User{
			Email:    req.Email,
			Password: req.Password,
		}

		if err := s.store.User().Create(u); err != nil {
			// User send incorrect data - 422 error
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		// hide password and render user without it
		u.Sanitize()
		// since user `u` is passed to `respond` method, need to set JSON tags in base User struct
		s.respond(w, r, http.StatusCreated, u)

	}
}

func (s *server) handleSessionsCreate() http.HandlerFunc {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		// find user by email and check that email is OK and passed password corresponds to encrypted one in store
		u, err := s.store.User().FindByEmail(req.Email)
		if err != nil || !u.ComparePasswords(req.Password) {
			s.error(w, r, http.StatusUnauthorized, errIncorrectEmailOrPassword)
			return
		}
		// return cookie to user after successful authentication
		// using gorilla/sessions package for that
		session, err := s.sessionStore.Get(r, sessionName)
		if err != nil {
			// return internal server error because problem is on our side
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		// Need to add middleware that gets user_id from session during every request, then goes to the store with user_id,
		// then gets user. If user exists - add it to context of current request. If not, return 401 error.

		// Load session if it has user ID. Otherwise, return unauthorized error or smth else
		session.Values["user_id"] = u.ID
		// Saving current session
		if err := s.sessionStore.Save(r, w, session); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, nil)
	}
}

// error is helper method to render any errors during work of handlers
// it will use another helper named `respond`
func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"error": err.Error()})
}

// respond is used for rendering of everything. `data` can have any type - set empty interface
func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
