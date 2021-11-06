package storage

import (
	"github.com/jinzhu/gorm"
	"github.com/super-radmir/menuxd_api/pkg/click"
	"github.com/super-radmir/menuxd_api/pkg/dish"
)

// DishStorage storage to the dish model.
type DishStorage struct {
	session *Session
	db      *gorm.DB
}

// setContext initialize the context to DishStorage.
func (s *DishStorage) setContext() {
	s.session = NewSession()
	s.db = s.session.Client
}

// Create create a new dish.
func (s DishStorage) Create(d *dish.Dish) error {
	s.setContext()

	if d.Name == "" || d.Pictures[0] == "" || d.Price < 0 {
		return ErrRequiredField
	}

	d.PicturesString = dish.SetString(d.Pictures)

	d.Available = true

	err := s.db.Create(d).Error
	if err != nil {
		return ErrNotInsert
	}

	for _, i := range d.Ingredients {
		i.DishID = d.ID
		err = s.db.Create(&i).Error
		if err != nil {
			return ErrNotInsert
		}
	}

	return nil
}

// CreateMany create multiple dishes to a client.
func (s DishStorage) CreateMany(clientID uint, d dish.Dishes) error {
	s.setContext()

	dishes := d.SetClientID(clientID)
	for _, nd := range dishes {
		nd.PicturesString = dish.SetString(nd.Pictures)
		err := s.db.Create(&nd).Error
		if err != nil {
			return ErrNotInsert
		}

		for _, ni := range nd.Ingredients {
			ni.DishID = nd.ID
			ni.ID = 0
			err := s.db.Create(&ni).Error
			if err != nil {
				return ErrNotInsert
			}
		}
	}

	return nil
}

// Update update a dish by ID.
func (s DishStorage) Update(id uint, updates map[string]interface{}) error {
	s.setContext()
	delete(updates, "client_id")

	iPictures, ok := updates["pictures"]
	if ok {
		i, ok := iPictures.([]interface{})
		if !ok {
			delete(updates, "PicturesString")
			delete(updates, "pictures")
		} else {
			pictures := []string{}
			for _, p := range i {
				pictures = append(pictures, p.(string))
			}
			updates["PicturesString"] = dish.SetString(pictures)
		}
	}

	err := s.db.Model(&dish.Dish{}).Where("id = ?", id).Updates(updates).Error
	if err != nil {
		return ErrNotUpdate
	}

	ingredients, ok := updates["ingredients"]
	if ok {
		ings, ok := ingredients.([]interface{})
		if ok {
			s.db.Delete(&dish.Ingredient{}, "dish_id = ?", id)
			for _, i := range ings {
				ing := dish.Ingredient{}
				newIng := i.(map[string]interface{})
				ing.Active = newIng["active"].(bool)
				ing.Name = newIng["name"].(string)
				ing.Price = newIng["price"].(float64)
				ing.DishID = id
				s.db.Create(&ing)
			}
		}
	}

	return nil
}

// Delete remove a dish by ID.
func (s DishStorage) Delete(id uint) error {
	s.setContext()

	err := s.db.Delete(&dish.Dish{}, "id = ?", id).Error
	if err != nil {
		return ErrNotDelete
	}

	return nil
}

// GetAll returns all stored dishes.
func (s DishStorage) GetAll(clientID uint) ([]dish.Dish, error) {
	s.setContext()

	dishes := []dish.Dish{}
	err := s.db.Model(&dish.Dish{}).Where("client_id = ?", clientID).
		Order("name").Find(&dishes).Error
	if err != nil {
		return []dish.Dish{}, ErrNotFound
	}

	for i := 0; i < len(dishes); i++ {
		dishes[i].Pictures = dish.SetSlice(dishes[i].PicturesString)
		s.db.Model(&dishes[i]).Related(&dishes[i].Ingredients)
	}

	return dishes, nil
}

// GetAllBackup returns all stored dishes.
func (s DishStorage) GetAllBackup(clientID uint) ([]dish.BaseDish, error) {
	s.setContext()

	dishes := []dish.Dish{}
	err := s.db.Model(&dish.Dish{}).Where("client_id = ?", clientID).
		Order("name").Find(&dishes).Error
	if err != nil {
		return []dish.BaseDish{}, ErrNotFound
	}

	result := []dish.BaseDish{}
	bd := dish.BaseDish{}
	for i := 0; i < len(dishes); i++ {
		dishes[i].Pictures = dish.SetSlice(dishes[i].PicturesString)
		s.db.Model(&dishes[i]).Related(&dishes[i].Ingredients)

		bd.Name = dishes[i].Name
		bd.Description = dishes[i].Description
		bd.Available = dishes[i].Available
		bd.Price = dishes[i].Price
		bd.Pictures = dishes[i].Pictures
		bd.Ingredients = dishes[i].Ingredients
		bd.Suggested = dishes[i].Suggested
		bd.IsHalf = dishes[i].IsHalf
		bd.HalfPrice = dishes[i].HalfPrice

		result = append(result, bd)
	}

	return result, nil
}

// GetAllWithPagination returns all stored dishes.
func (s DishStorage) GetAllWithPagination(clientID uint, page int64) (dish.Dishes, int, error) {
	s.setContext()

	var err error
	var limit int64 = 12
	dishes := dish.Dishes{}

	err = s.db.Limit(limit).Offset((page-1)*limit).Order("name").Find(
		&dishes,
		"client_id = ?",
		clientID,
	).Error
	if err != nil {
		return dish.Dishes{}, 0, ErrNotFound
	}

	var cs CategoryStorage
	for i := 0; i < len(dishes); i++ {
		dishes[i].Pictures = dish.SetSlice(dishes[i].PicturesString)
		s.db.Model(&dishes[i]).Related(&dishes[i].Ingredients)
		storedCategory, err := cs.GetByID(dishes[i].CategoryID)
		if err != nil {
			continue
		}
		dishes[i].Category = &storedCategory
	}

	var total int
	err = s.db.Model(&dish.Dish{}).Where("client_id = ?", clientID).
		Count(&total).Error
	if err != nil {
		return dishes, 0, ErrNotFound
	}

	return dishes, total, nil
}

// GetAllByCategory returns dishes by Category ID.
func (s DishStorage) GetAllByCategory(categoryID uint) (dish.Dishes, error) {
	s.setContext()

	dishes := dish.Dishes{}
	err := s.db.Order("name").Find(&dishes, "category_id = ?", categoryID).Error
	if err != nil {
		return []dish.Dish{}, ErrNotFound
	}

	var cs CategoryStorage
	for i := 0; i < len(dishes); i++ {
		dishes[i].Pictures = dish.SetSlice(dishes[i].PicturesString)
		s.db.Model(&dishes[i]).Related(&dishes[i].Ingredients)
		storedCategory, err := cs.GetByID(dishes[i].CategoryID)
		if err != nil {
			continue
		}
		dishes[i].Category = &storedCategory
	}

	return dishes, nil
}

// GetAllActiveByCategory returns active dishes by Category ID.
func (s DishStorage) GetAllActiveByCategory(categoryID uint) (dish.Dishes, error) {
	s.setContext()

	dishes := dish.Dishes{}
	err := s.db.Order("name").Where("available = ?", true).Find(&dishes, "category_id = ?", categoryID).Error
	if err != nil {
		return []dish.Dish{}, ErrNotFound
	}

	var cs CategoryStorage
	for i := 0; i < len(dishes); i++ {
		dishes[i].Pictures = dish.SetSlice(dishes[i].PicturesString)
		s.db.Model(&dishes[i]).Related(&dishes[i].Ingredients)
		storedCategory, err := cs.GetByID(dishes[i].CategoryID)
		if err != nil {
			continue
		}
		dishes[i].Category = &storedCategory
	}

	return dishes, nil
}

// GetByID returns a dish by ID.
func (s DishStorage) GetByID(id uint) (dish.Dish, error) {
	s.setContext()

	d := dish.Dish{}
	err := s.db.First(&d, "id = ?", id).Error
	if err != nil {
		return dish.Dish{}, ErrNotFound

	}
	var cs CategoryStorage
	d.Pictures = dish.SetSlice(d.PicturesString)
	d.PicturesString = ""
	storedCategory, err := cs.GetByID(d.CategoryID)
	d.Category = &storedCategory
	s.db.Model(&d).Related(&d.Ingredients)

	return d, nil
}

// GetSuggested returns suggested drinks.
func (s DishStorage) GetSuggested(categoryID uint) (dish.Dishes, error) {
	s.setContext()

	dishes := dish.Dishes{}
	err := s.db.Where("suggested = ?", true).Order("name").
		Find(&dishes, "category_id = ?", categoryID).Error
	if err != nil {
		return []dish.Dish{}, ErrNotFound
	}

	var cs CategoryStorage
	for i := 0; i < len(dishes); i++ {
		dishes[i].Pictures = dish.SetSlice(dishes[i].PicturesString)
		s.db.Model(&dishes[i]).Related(&dishes[i].Ingredients)
		storedCategory, err := cs.GetByID(dishes[i].CategoryID)
		if err != nil {
			continue
		}
		dishes[i].Category = &storedCategory
	}

	return dishes, nil
}

// AddClick create a new click.
func (s DishStorage) AddClick(suggestedID uint) error {
	s.setContext()

	c := click.Click{}
	c.TypeID = suggestedID
	c.Type = click.Suggested

	err := s.db.Create(&c).Error
	if err != nil {
		return ErrNotInsert
	}

	return nil
}

// GetClicks get all clicks by client id.
func (s DishStorage) GetClicks(clientID uint) ([]click.Click, error) {
	s.setContext()

	clicks := []click.Click{}

	err := s.db.Model(&click.Click{}).
	Where("type = ?", click.Suggested).
	Find(&clicks, "type_id = ?", clientID).Error
	if err != nil {
		return []click.Click{}, ErrNotFound
	}

	return clicks, nil
}
