package waiter

import (
	"regexp"

	"github.com/super-radmir/menuxd_api/pkg/model"
)

// Storage handle the CRUD operations with Waiters.
type Storage interface {
	Create(waiter *Waiter) error
	Update(id uint, waiter *Waiter) error
	Delete(id uint) error
	GetAll(clientID uint) (Waiters, error)
	GetByID(id uint) (Waiter, error)
}

type Waiter struct {
	model.Model
	Name     string `bson:"name" json:"name"`
	PIN      string `gorm:"default:'1234'" bson:"pin" json:"pin"`
	ClientID uint   `bson:"client_id" json:"client_id"`
}

// VerifyPIN validate a PIN, like 1234
func (w *Waiter) VerifyPIN() bool {
	re, _ := regexp.Compile("[0-9]{4}")
	return re.MatchString(w.PIN) && len(w.PIN) == 4
}

// Waiters alias for a slice of Waiters.
type Waiters []Waiter

// New returns a instance of Waiter with default configuration.
func New() *Waiter {
	return &Waiter{}
}
