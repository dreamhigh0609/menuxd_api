package storage

import (
	"github.com/jinzhu/gorm"

	"github.com/super-radmir/menuxd_api/pkg/bill"
)

// BillStorage storage to the bill model
type BillStorage struct {
	session *Session
	db      *gorm.DB
}

// setContext initialize the context to BillStorage
func (s *BillStorage) setContext() {
	s.session = NewSession()
	s.db = s.session.Client
}

// Create create a new bill
func (s BillStorage) Create(b *bill.Bill) error {
	s.setContext()

	err := s.db.Create(b).Error
	if err != nil {
		return ErrNotInsert
	}

	return nil
}

// Update update bill by ID
func (s BillStorage) Update(id uint, b *bill.Bill) error {
	s.setContext()

	updates := map[string]interface{}{
		"value": b.Value,
		"paid":  b.Paid,
	}

	err := s.db.Model(&bill.Bill{}).Where("id = ?", id).Updates(updates).Error
	if err != nil {
		return ErrNotUpdate
	}

	return nil
}

// Delete remove a bill by ID
func (s BillStorage) Delete(id uint) error {
	s.setContext()

	err := s.db.Delete(&bill.Bill{}, "id ").Error
	if err != nil {
		return ErrNotDelete
	}

	return nil
}

// GetAll returns all stored bills
func (s BillStorage) GetAll(clientId uint) (bill.Bills, error) {
	s.setContext()

	bills := bill.Bills{}
	err := s.db.Find(&bills, "client_id = ?", clientId).Error
	if err != nil {
		return []bill.Bill{}, ErrNotFound
	}

	return bills, nil
}

// GetByID returns a bill by ID
func (s BillStorage) GetByID(id uint) (bill.Bill, error) {
	s.setContext()

	b := bill.Bill{}
	err := s.db.First(&b, "id = ?", id).Error
	if err != nil {
		return bill.Bill{}, ErrNotFound
	}

	return b, nil
}
