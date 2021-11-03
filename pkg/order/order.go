package order

import (
	"errors"

	"github.com/super-radmir/menuxd_api/pkg/dish"
	"github.com/super-radmir/menuxd_api/pkg/model"
	"github.com/super-radmir/menuxd_api/pkg/table"
)

// Errors.
var (
	ErrParseFailer = errors.New("parse orders failed")
	ErrInvalidType = errors.New("invalid type")
)

// Storage handle Order's CRUD.
type Storage interface {
	Create(o *Order) (Order, error)
	Add(id uint, items []Item) error
	Update(id uint, o *Order) error
	GetAll(clientID uint) ([]Order, error)
	GetAllActive(clientID uint) ([][]Order, error)
	GetByID(id uint) (Order, error)
	PatchItem(id uint, updates map[string]interface{}) error
}

// Order is a Client's request.
type Order struct {
	model.Model
	ClientID uint         `json:"client_id"`
	TableID  uint         `json:"table_id"`
	Table    *table.Table `json:"table"`
	Canceled bool         `json:"canceled"`
	Items    []Item       `json:"items"`
}

// Item is a element to order.
type Item struct {
	model.Model
	OrderID             uint                 `json:"order_id"`
	Mount               uint                 `json:"mount"`
	Active              bool                 `gorm:"default:true" json:"active"`
	Ready               bool                 `gorm:"default:true" json:"ready"`
	DishID              uint                 `json:"dish_id"`
	Ingredients         []dish.Ingredient    `json:"ingredients"`
	SelectedIngredients []IngredientSelected `json:"selected_ingredients"`
	Dish                *dish.Dish           `json:"dish,omitempty"`
	Takeaway            bool                 `json:"takeaway"`
	Locked              bool                 `gorm:"-" json:"locked,omitempty"`
}

// IngredientSelected is an ingredient selected in an order.
type IngredientSelected struct {
	model.Model
	ItemID       uint             `json:"item_id"`
	Active       bool             `json:"active"`
	IngredientID uint             `json:"ingredient_id"`
	Ingredient   *dish.Ingredient `json:"ingredient,omitempty"`
}

// Orders is an alias to a slice of order.
type Orders []Order
