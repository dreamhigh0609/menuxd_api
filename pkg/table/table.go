package table

import "github.com/super-radmir/menuxd_api/pkg/model"

// Storage handle the CRUD operations with Tables.
type Storage interface {
	Create(table *Table) error
	Update(id uint, table *Table) error
	Delete(id uint) error
	GetAll(clientID uint) (Tables, error)
	GetByID(id uint) (Table, error)
}

// storage is a instance of Storage interface.
var storage Storage

// SetStorage set a new storage.
func SetStorage(s Storage) {
	storage = s
}

// Table represents a table or bar in a restaurant.
type Table struct {
	model.Model
	Number      uint   `bson:"number" json:"number"`
	Type        string `gorm:"default:'table'" bson:"type" json:"type"`
	Available   bool   `gorm:"default:true" bson:"available" json:"available"`
	ClientID    uint   `bson:"client_id" json:"client_id"`
	CallsWaiter bool   `bson:"calls_waiter" json:"calls_waiter"`
	AsksForBill bool   `bson:"asks_for_bill" json:"asks_for_bill"`
}

// Tables alias for a slice of Tables.
type Tables []Table

// New returns a instance of Table with default configuration.
func New() *Table {
	return &Table{}
}
