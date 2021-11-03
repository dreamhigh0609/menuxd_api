package storage

import (
	"github.com/jinzhu/gorm"

	"github.com/super-radmir/menuxd_api/pkg/category"
)

// CategoryStorage storage to the category model.
type CategoryStorage struct {
	db      *gorm.DB
	session *Session
}

// setContext initialize the context to CategoryStorage.
func (s *CategoryStorage) setContext() {
	s.session = NewSession()
	s.db = s.session.Client
}

// Create create a new category.
func (s CategoryStorage) Create(c *category.Category) error {
	s.setContext()

	if c.Title == "" || c.Picture == "" {
		return ErrRequiredField
	}

	c.Active = true
	err := s.db.Create(c).Error
	if err != nil {
		return ErrNotInsert
	}

	return nil
}

// CreateMany create multiple categories to a client.
func (s CategoryStorage) CreateMany(clientID uint, categories []category.Category) error {
	s.setContext()

	for i := 0; i < len(categories); i++ {
		categories[i].ClientID = clientID
	}

	for _, nc := range categories {
		err := s.db.Create(&nc).Error
		if err != nil {
			return ErrNotInsert
		}
	}

	return nil
}

// Update update category by ID.
func (s CategoryStorage) Update(id uint, c *category.Category) error {
	s.setContext()

	if c.Title == "" || c.Picture == "" {
		return ErrRequiredField
	}

	updates := map[string]interface{}{
		"title":      c.Title,
		"picture":    c.Picture,
		"suggested1": c.Suggested1,
		"suggested2": c.Suggested2,
		"suggested3": c.Suggested3,
		"priority":   c.Priority,
	}

	err := s.db.Model(&category.Category{}).Where("id = ?", id).
		Updates(updates).Error
	if err != nil {
		return ErrNotUpdate
	}

	return nil
}

// Patch update part of the category by ID.
func (s CategoryStorage) Patch(id uint, updates map[string]interface{}) error {
	s.setContext()

	delete(updates, "client_id")

	err := s.db.Model(&category.Category{}).Where("id = ?", id).
		Updates(updates).Error
	if err != nil {
		return ErrNotUpdate
	}

	return nil
}

// UpdatePositions update positions to categories.
func (s CategoryStorage) UpdatePositions(categories []category.Category) error {
	s.setContext()

	for _, c := range categories {
		err := s.db.Model(&category.Category{}).Where("id = ?", c.ID).
			Update("position", c.Position).Error
		if err != nil {
			return ErrNotUpdate
		}
	}

	return nil
}

// Delete remove a category by ID.
func (s CategoryStorage) Delete(id uint) error {
	s.setContext()

	err := s.db.Delete(&category.Category{}, "id = ?", id).Error
	if err != nil {
		return ErrNotDelete
	}

	return nil
}

// GetAll returns all stored categories.
func (s CategoryStorage) GetAll(clientID uint) (category.Categories, error) {
	s.setContext()

	categories := category.Categories{}
	err := s.db.Order("position ASC").Order("title").
		Find(&categories, "client_id = ?", clientID).Error
	if err != nil {
		return category.Categories{}, ErrNotFound
	}
	return categories, nil
}

// GetAllActive returns all active categories.
func (s CategoryStorage) GetAllActive(clientID uint) (category.Categories, error) {
	s.setContext()

	categories := category.Categories{}
	err := s.db.Order("position ASC").Order("title").Where("active = ?", true).
		Find(&categories, "client_id = ?", clientID).Error
	if err != nil {
		return category.Categories{}, ErrNotFound
	}
	return categories, nil
}

// GetAllBackup returns all stored dishes.
func (s CategoryStorage) GetAllBackup(clientID uint) (
	[]category.BaseCategory, error,
) {
	s.setContext()

	categories := []category.BaseCategory{}
	err := s.db.Model(&category.Category{}).Select(
		"title, priority, active, picture, suggested1, suggested2, suggested3, position",
	).Where("client_id = ?", clientID).Scan(&categories).Error
	if err != nil {
		return []category.BaseCategory{}, ErrNotFound
	}

	return categories, nil
}

// GetByID returns a category by ID.
func (s CategoryStorage) GetByID(id uint) (category.Category, error) {
	s.setContext()

	c := category.Category{}
	err := s.db.First(&c, "id = ?", id).Error
	if err != nil {
		return category.Category{}, ErrNotFound
	}

	return c, nil
}
