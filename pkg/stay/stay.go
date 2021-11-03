package stay

import (
	"errors"

	"github.com/super-radmir/menuxd_api/pkg/model"
)

// Storage handle the CRUD operations with Saty.
type Storage interface {
	Create(stay *Stay) error
	GetAll() ([]Stay, error)
	GetByClient(clientID uint) ([]Stay, error)
	GetByID(id uint) (Stay, error)
}

// Stay is the customers stay time in the app.
type Stay struct {
	model.Model
	Time     float64 `bson:"time" json:"time,omitempty"`
	ClientID uint    `bson:"client_id" json:"client_id,omitempty"`
}

// IsValid verifies that the time is valid.
func (s Stay) IsValid() bool {
	return s.Time >= 0
}

// BeforeSave checks that the time is valid before inserting.
func (s *Stay) BeforeSave() (err error) {
	if !s.IsValid() {
		err = errors.New("can't save invalid data")
	}

	return
}
