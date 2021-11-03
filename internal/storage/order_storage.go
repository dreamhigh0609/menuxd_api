package storage

import (
	"github.com/jinzhu/gorm"
	"github.com/super-radmir/menuxd_api/pkg/dish"
	"github.com/super-radmir/menuxd_api/pkg/order"
	"github.com/super-radmir/menuxd_api/pkg/table"
)

// OrderStorage storage to the dish model.
type OrderStorage struct {
	session *Session
	db      *gorm.DB
}

// setContext initialize the context to OrderStorage.
func (s *OrderStorage) setContext() {
	s.session = NewSession()
	s.db = s.session.Client
}

// Create create a new order.
func (s OrderStorage) Create(o *order.Order) (order.Order, error) {
	s.setContext()

	o.Items = []order.Item{}
	if o.Table != nil {
		o.TableID = o.Table.ID
	}
	o.Table = nil

	err := s.db.Create(o).Error
	if err != nil {
		return order.Order{}, ErrNotInsert
	}

	return *o, nil
}

// Add update a dish by ID.
func (s OrderStorage) Add(id uint, items []order.Item) error {
	s.setContext()

	for _, i := range items {
		if i.Dish != nil {
			i.DishID = i.Dish.ID
		}
		ingredients := []dish.Ingredient{}
		i.Dish = nil
		i.OrderID = id
		i.Active = true
		ingredients = i.Ingredients[:]

		i.Ingredients = nil
		i.SelectedIngredients = nil
		err := s.db.Create(&i).Error
		if err != nil {
			return ErrNotInsert
		}

		for _, ing := range ingredients {
			is := order.IngredientSelected{}
			is.ItemID = i.ID
			is.Active = ing.Active
			is.IngredientID = ing.ID

			err = s.db.Create(&is).Error
			if err != nil {
				return ErrNotInsert
			}
		}
	}

	return nil
}

// PatchItem set item's status.
func (s OrderStorage) PatchItem(id uint, updates map[string]interface{}) error {
	s.setContext()

	delete(updates, "order_id")
	delete(updates, "dish_id")
	delete(updates, "takeaway")
	delete(updates, "mount")

	err := s.db.Model(&order.Item{}).Where("id = ?", id).
		Updates(updates).Error
	if err != nil {
		return ErrNotUpdate
	}

	return nil
}

// Update set order's canceled.
func (s OrderStorage) Update(id uint, o *order.Order) error {
	s.setContext()

	o.ID = id
	err := s.db.Model(o).Update("canceled", o.Canceled).Error
	if err != nil {
		return ErrNotUpdate
	}

	return nil
}

func (s OrderStorage) getIngredientsByItem(itemID uint) ([]order.IngredientSelected, error) {
	ingredients := []order.IngredientSelected{}
	err := s.db.Model(&order.IngredientSelected{}).
		Find(&ingredients, "item_id = ?", itemID).Error

	if err != nil {
		return []order.IngredientSelected{}, ErrNotInsert
	}

	result := []order.IngredientSelected{}
	for _, ing := range ingredients {
		storedIngredient := dish.Ingredient{}
		err := s.db.Model(&dish.Ingredient{}).
			First(&storedIngredient, "id = ?", ing.IngredientID).Error
		if err != nil {
			continue
		}

		ing.Ingredient = &storedIngredient

		result = append(result, ing)
	}

	return result, nil
}

// GetAll returns all stored orders.
func (s OrderStorage) GetAll(clientID uint) ([]order.Order, error) {
	s.setContext()

	orders := []order.Order{}
	err := s.db.Model(&order.Order{}).Order("created_at").
		Find(&orders, "client_id = ?", clientID).Error
	if err != nil {
		return []order.Order{}, ErrNotFound
	}

	result := []order.Order{}
	var ds DishStorage

	for _, o := range orders {
		t := table.Table{}
		t.ID = o.TableID
		s.db.Model(&table.Table{}).First(&t, "id = ?", o.TableID)
		o.Table = &t
		s.db.Model(&order.Item{}).Find(&o.Items, "order_id = ?", o.ID)
		items := []order.Item{}
		for _, i := range o.Items {
			i.SelectedIngredients, _ = s.getIngredientsByItem(i.ID)
			storedDish, err := ds.GetByID(i.DishID)
			if err != nil {
				i.Dish = &dish.Dish{}
			} else {
				i.Dish = &storedDish
			}

			items = append(items, i)
		}

		o.Items = items
		result = append(result, o)
	}

	return result, nil
}

// GetAllActive returns active orders.
func (s OrderStorage) GetAllActive(clientID uint) ([][]order.Order, error) {
	s.setContext()

	orders := []order.Order{}
	err := s.db.Model(&order.Order{}).Where("canceled = false").
		Order("id ASC").
		Find(&orders, "client_id = ?", clientID).Error
	if err != nil {
		return [][]order.Order{}, ErrNotFound
	}

	result := []order.Order{}
	var ds DishStorage

	for _, o := range orders {
		t := table.Table{}
		t.ID = o.TableID
		s.db.Model(&table.Table{}).First(&t, "id = ?", o.TableID)
		o.Table = &t
		s.db.Model(&order.Item{}).Order("id ASC").Find(&o.Items, "order_id = ?", o.ID)
		items := []order.Item{}
		for _, i := range o.Items {
			i.SelectedIngredients, _ = s.getIngredientsByItem(i.ID)
			storedDish, err := ds.GetByID(i.DishID)
			if err != nil {
				i.Dish = &dish.Dish{}
			} else {
				i.Dish = &storedDish
			}

			items = append(items, i)
		}

		o.Items = items
		result = append(result, o)
	}

	m := make(map[uint][]order.Order)
	for _, r := range result {
		m[r.TableID] = append(m[r.TableID], r)
	}

	tables := [][]order.Order{}
	for _, t := range m {
		tables = append(tables, t)
	}

	return tables, nil
}

// GetByID returns a dish by ID.
func (s OrderStorage) GetByID(id uint) (order.Order, error) {
	s.setContext()

	o := order.Order{}
	err := s.db.First(&o, "id = ?", id).Error
	if err != nil {
		return order.Order{}, ErrNotFound
	}

	err = s.db.Model(&order.Item{}).Find(&o.Items, "order_id = ?", o.ID).Error
	if err != nil {
		return order.Order{}, ErrNotFound
	}

	var ds DishStorage
	items := []order.Item{}
	for _, i := range o.Items {
		i.SelectedIngredients, _ = s.getIngredientsByItem(i.ID)
		storedDish, err := ds.GetByID(i.DishID)
		if err != nil {
			i.Dish = &dish.Dish{}
		} else {
			i.Dish = &storedDish
		}

		items = append(items, i)
	}

	o.Items = items
	t := table.Table{}

	err = s.db.Model(&table.Table{}).First(&t, "id = ?", o.TableID).Error
	if err != nil {
		return order.Order{}, ErrNotFound
	}

	o.Table = &t

	return o, nil
}
