package storage

import (
	"github.com/jinzhu/gorm"

	"github.com/super-radmir/menuxd_api/pkg/waiter"
)

// WaiterStorage storage to the waiter model
type WaiterStorage struct {
	session *Session
	db      *gorm.DB
}

// setContext initialize the context to WaiterStorage.
func (s *WaiterStorage) setContext() {
	s.session = NewSession()
	s.db = s.session.Client
}

// Create create a new waiter.
func (s WaiterStorage) Create(w *waiter.Waiter) error {
	s.setContext()

	if w.Name == "" {
		return ErrRequiredField
	}

	if !w.VerifyPIN() {
		return ErrInvalidPIN
	}

	err := s.db.Create(w).Error
	if err != nil {
		return ErrNotInsert
	}
	return nil
}

// Update update a waiter by ID.
func (s WaiterStorage) Update(id uint, w *waiter.Waiter) error {
	s.setContext()

	if !w.VerifyPIN() {
		return ErrInvalidPIN
	}

	updates := map[string]interface{}{
		"name": w.Name,
		"pin":  w.PIN,
	}

	w.ID = id
	err := s.db.Model(w).Updates(updates).Error
	if err != nil {
		return ErrNotUpdate
	}

	return nil
}

// Delete remove a waiter by ID.
func (s WaiterStorage) Delete(id uint) error {
	s.setContext()

	err := s.db.Delete(&waiter.Waiter{}, "id = ?", id).Error
	if err != nil {
		return ErrNotDelete
	}

	return nil
}

// GetAll returns all stored waiters.
func (s WaiterStorage) GetAll(clientID uint) (waiter.Waiters, error) {
	s.setContext()

	waiters := waiter.Waiters{}
	err := s.db.Order("name").Find(&waiters, "client_id = ?", clientID).Error
	if err != nil {
		return []waiter.Waiter{}, ErrNotFound
	}

	return waiters, nil
}

// GetByID returns a waiter by ID.
func (s WaiterStorage) GetByID(id uint) (waiter.Waiter, error) {
	s.setContext()

	w := waiter.Waiter{}
	err := s.db.First(&w, "id = ?", id).Error
	if err != nil {
		return waiter.Waiter{}, ErrNotFound
	}

	return w, nil
}
