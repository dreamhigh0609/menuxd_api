package category

import (
	"errors"

	"github.com/super-radmir/menuxd_api/pkg/model"
)

// Storage handle the CRUD operations with Categories.
type Storage interface {
	Create(category *Category) error
	Update(id uint, category *Category) error
	CreateMany(clientID uint, categories []Category) error
	Delete(id uint) error
	GetAll(clientID uint) (Categories, error)
	GetAllActive(clientID uint) (Categories, error)
	GetByID(id uint) (Category, error)
	Patch(id uint, updates map[string]interface{}) error
	GetAllBackup(clientID uint) ([]BaseCategory, error)
	UpdatePositions(categories []Category) error
}

// BaseCategory is a lite category for a dish.
type BaseCategory struct {
	Title      string `bson:"title" json:"title"`
	Picture    string `bson:"picture" json:"picture"`
	Active     bool   `gorm:"default:true" json:"active"`
	Priority   bool   `gorm:"default:false" json:"priority"`
	Suggested1 *uint  `gorm:"default:null" bson:"suggested1,omitempty" json:"suggested1,omitempty"`
	Suggested2 *uint  `gorm:"default:null" bson:"suggested2,omitempty" json:"suggested2,omitempty"`
	Suggested3 *uint  `gorm:"default:null" bson:"suggested3,omitempty" json:"suggested3,omitempty"`
	Position   uint   `gorm:"default:1" bson:"position" json:"position,omitempty"`
}

// Category for a dish.
type Category struct {
	model.Model
	BaseCategory
	ClientID uint `bson:"client_id" json:"client_id"`
}

// IsValid checks that the suggested IDs are unique.
func (c Category) IsValid() bool {
	var suggested1 uint
	var suggested2 uint
	var suggested3 uint

	if c.Suggested1 != nil {
		suggested1 = *c.Suggested1
	}

	if c.Suggested2 != nil {
		suggested2 = *c.Suggested2
	}

	if c.Suggested3 != nil {
		suggested3 = *c.Suggested3
	}

	if suggested1 > 0 {
		return suggested1 != suggested2 && suggested1 != suggested3
	}

	if suggested2 > 0 {
		return suggested2 != suggested3
	}

	return true
}

// BeforeSave checks that the category is valid before inserting.
func (c *Category) BeforeSave() (err error) {
	if !c.IsValid() {
		err = errors.New("can't save invalid data")
	}
	return
}

// Categories alias for a slice of Categories.
type Categories []Category

// New returns a instance of Bill with default configuration.
func New() *Category {
	return &Category{}
}
