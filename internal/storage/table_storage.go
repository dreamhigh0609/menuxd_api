package storage

import (
	"github.com/jinzhu/gorm"

	"github.com/super-radmir/menuxd_api/pkg/table"
)

// TableStorage storage to the table model
type TableStorage struct {
	session *Session
	db      *gorm.DB
}

// setContext initialize the context to TableStorage
func (s *TableStorage) setContext() {
	s.session = NewSession()
	s.db = s.session.Client
}

// Create create a new table
func (s TableStorage) Create(t *table.Table) error {
	s.setContext()

	t.Available = true

	err := s.db.Create(t).Error
	if err != nil {
		return ErrNotInsert
	}
	return nil
}

// Update update a table by ID
func (s TableStorage) Update(id uint, t *table.Table) error {
	s.setContext()

	t.ID = id

	err := s.db.Model(t).Update("available", t.Available).Error
	if err != nil {
		return ErrNotUpdate
	}

	return nil
}

// Delete remove a table by ID
func (s TableStorage) Delete(id uint) error {
	s.setContext()

	err := s.db.Delete(&table.Table{}, "id = ?", id).Error
	if err != nil {
		return ErrNotDelete
	}

	return nil
}

// GetAll returns all stored tables
func (s TableStorage) GetAll(clientID uint) (table.Tables, error) {
	s.setContext()

	tables := table.Tables{}
	err := s.db.Order("type DESC").Order("number").Find(&tables, "client_id = ?", clientID).Error
	if err != nil {
		return []table.Table{}, ErrNotFound
	}

	return tables, nil
}

// GetByID returns a table by ID
func (s TableStorage) GetByID(id uint) (table.Table, error) {
	s.setContext()

	t := table.Table{}
	err := s.db.First(&t, "id = ?", id).Error
	if err != nil {
		return table.Table{}, ErrNotFound
	}

	return t, nil
}
