package store

// Store is an interface for store
type Store interface {
	User() UserRepository
}
