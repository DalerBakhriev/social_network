package sqlstore

import (
	"database/sql"

	"github.com/DalerBakhriev/social_network/internal/app/store"
	_ "github.com/go-sql-driver/mysql" // driver import
)

// Store ..
type Store struct {
	db             *sql.DB
	userRepository *UserRepository
}

// New ...
func New(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

// User returns user repository to work with sql store
func (s *Store) User() store.UserRepository {

	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = &UserRepository{
		store: s,
	}

	return s.userRepository
}
