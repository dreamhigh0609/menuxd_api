package bill

import (
	"github.com/super-radmir/menuxd_api/pkg/model"
	"github.com/super-radmir/menuxd_api/pkg/order"
)

// Storage handle the CRUD operations with Bills.
type Storage interface {
	Create(bill *Bill) error
	Update(id uint, bill *Bill) error
	Delete(id uint) error
	GetAll(clientID uint) (Bills, error)
	GetByID(id uint) (Bill, error)
}

type Bill struct {
	model.Model
	Value    uint          `bson:"value" json:"value"`
	Paid     bool          `bson:"paid" json:"paid"`
	Orders   []order.Order `bson:"orders" json:"orders"`
	TableID  uint          `bson:"table_id" json:"table_id"`
	ClientID uint          `bson:"client_id" json:"client_id"`
}

// Bills alias for a slice of Bills.
type Bills []Bill

// New returns a instance of Bill with default configuration.
func New() *Bill {
	return &Bill{}
}
