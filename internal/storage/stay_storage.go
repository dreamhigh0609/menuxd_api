package storage

import (
	"github.com/jinzhu/gorm"

	"github.com/super-radmir/menuxd_api/pkg/stay"
)

// StayStorage storage to the stay model.
type StayStorage struct {
	session *Session
	db      *gorm.DB
}

// setContext initialize the context to StayStorage.
func (s *StayStorage) setContext() {
	s.session = NewSession()
	s.db = s.session.Client
}

// Create create a new stay.
func (s StayStorage) Create(st *stay.Stay) error {
	s.setContext()

	err := s.db.Create(st).Error
	if err != nil {
		return ErrNotInsert
	}
	return nil
}

// GetAll returns all stored stay.
func (s StayStorage) GetAll() ([]stay.Stay, error) {
	s.setContext()

	result := []stay.Stay{}
	err := s.db.Order("created_at DESC").Find(&result).Error
	if err != nil {
		return []stay.Stay{}, ErrNotFound
	}

	return result, nil
}

// GetByClient returns all stored stay by client ID.
func (s StayStorage) GetByClient(clientID uint) ([]stay.Stay, error) {
	s.setContext()

	result := []stay.Stay{}
	err := s.db.Order("created_at DESC").Find(&result, "client_id = ?", clientID).Error
	if err != nil {
		return []stay.Stay{}, ErrNotFound
	}

	return result, nil
}

// GetByID returns a stay by ID.
func (s StayStorage) GetByID(id uint) (stay.Stay, error) {
	s.setContext()

	st := stay.Stay{}
	err := s.db.First(&st, "id = ?", id).Error
	if err != nil {
		return stay.Stay{}, ErrNotFound
	}

	return st, nil
}
