package dish

import (
	"strings"

	"github.com/super-radmir/menuxd_api/pkg/category"
	"github.com/super-radmir/menuxd_api/pkg/click"
	"github.com/super-radmir/menuxd_api/pkg/model"
)

// Storage handle the CRUD operations with Dishes.
type Storage interface {
	Create(dish *Dish) error
	CreateMany(clientID uint, dishes Dishes) error
	Update(id uint, updates map[string]interface{}) error
	Delete(id uint) error
	GetAll(clientID uint) ([]Dish, error)
	GetAllBackup(clientID uint) ([]BaseDish, error)
	GetAllWithPagination(clientID uint, page int64) (Dishes, int, error)
	GetAllByCategory(categoryID uint) (Dishes, error)
	GetAllActiveByCategory(categoryID uint) (Dishes, error)
	GetByID(id uint) (Dish, error)
	GetSuggested(categoryID uint) (Dishes, error)
	AddClick(suggestedID uint) error
	GetClicks(clientID uint) ([]click.Click, error)
}

// BaseDish lite version of a dish.
type BaseDish struct {
	Name           string       `bson:"name" json:"name,omitempty"`
	Description    string       `bson:"description,omitempty" json:"description,omitempty"`
	Available      bool         `gorm:"default:true" bson:"available" json:"available"`
	Price          float64      `bson:"price" json:"price"`
	IsHalf         bool         `gorm:"default:false" bson:"is_half" json:"is_half"`
	HalfPrice      *float64     `gorm:"default:null" bson:"half_price,omitempty" json:"half_price,omitempty"`
	Pictures       []string     `gorm:"-" bson:"pictures" json:"pictures,omitempty"`
	PicturesString string       `gorm:"column:pictures" bson:"pictures" json:"pictures_string,omitempty"`
	Ingredients    []Ingredient `gorm:"-" json:"ingredients" json:"ingredients"`
	Suggested      bool         `gorm:"default:false" bson:"suggested" json:"suggested"`
}

// Ingredient to the dishes.
type Ingredient struct {
	model.Model
	DishID  uint    `json:"dish_id"`
	Name    string  `json:"name"`
	OrderID *uint   `json:"order_id"`
	Active  bool    `json:"active"`
	Price   float64 `json:"price"`
}

// Dish meal from a restaurant.
type Dish struct {
	model.Model
	BaseDish
	CategoryID uint               `bson:"category_id,omitempty" json:"category_id,omitempty"`
	Category   *category.Category `gorm:"-" bson:"category,omitempty" json:"category,omitempty"`
	ClientID   uint               `bson:"client_id" json:"client_id,omitempty"`
}

// SetSlice split strings into slices.
func SetSlice(arg string) []string {
	if arg != "" {
		return strings.Split(arg, ",")
	}

	return []string{"", "", ""}
}

// SetString split strings into slices.
func SetString(arg []string) string {
	if len(arg) >= 0 {
		return strings.Join(arg, ",")
	}

	return ""
}

// Dishes alias for a slice of Dishes.
type Dishes []Dish

// SetClientID set ClientID for all Dishes to the clientID given.
func (dishes *Dishes) SetClientID(clientID uint) (result []Dish) {
	for _, d := range *dishes {
		d.ClientID = clientID
		result = append(result, d)
	}
	return
}

// New returns a instance of Dish with default configuration.
func New() *Dish {
	return &Dish{}
}
