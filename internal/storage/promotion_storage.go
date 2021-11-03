package storage

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/super-radmir/menuxd_api/pkg/click"
	"github.com/super-radmir/menuxd_api/pkg/promotion"
)

// PromotionStorage storage to the promotion model
type PromotionStorage struct {
	session *Session
	db      *gorm.DB
}

// setContext initialize the context to PromotionStorage
func (s *PromotionStorage) setContext() {
	s.session = NewSession()
	s.db = s.session.Client
}

// Create create a new Promotion
func (s PromotionStorage) Create(p *promotion.Promotion) error {
	s.setContext()

	if p.Title == "" ||
		p.Picture == "" ||
		p.DishID == 0 ||
		p.StartAt == "" ||
		p.EndAt == "" {
		return ErrRequiredField
	}

	p.DaysString = promotion.SetDaysString(p.Days)
	err := s.db.Create(p).Error
	if err != nil {
		return ErrNotInsert
	}

	return nil
}

// Update update a promotion by ID.
func (s PromotionStorage) Update(id uint, p *promotion.Promotion) error {
	s.setContext()

	if p.Title == "" ||
		p.Picture == "" ||
		p.DishID == 0 ||
		p.StartAt == "" ||
		p.EndAt == "" {
		return ErrRequiredField
	}

	p.DaysString = promotion.SetDaysString(p.Days)

	updates := map[string]interface{}{
		"title":      p.Title,
		"picture":    p.Picture,
		"DaysString": p.DaysString,
		"dish_id":    p.DishID,
		"start_at":   p.StartAt,
		"end_at":     p.EndAt,
	}

	err := s.db.Model(&promotion.Promotion{}).Where("id = ?", id).Updates(updates).Error
	if err != nil {
		return ErrNotUpdate
	}

	return nil
}

// Delete remove a promotion by ID.
func (s PromotionStorage) Delete(id uint) error {
	s.setContext()

	err := s.db.Delete(&promotion.Promotion{}, "id = ?", id).Error
	if err != nil {
		return ErrNotDelete
	}

	return nil
}

// AddClick create a new click.
func (s PromotionStorage) AddClick(promotionID uint) error {
	s.setContext()

	c := click.Click{}
	c.TypeID = promotionID
	c.Type = click.Promotion

	err := s.db.Create(&c).Error
	if err != nil {
		return ErrNotInsert
	}

	return nil
}

// GetAll returns all stored promotions.
func (s PromotionStorage) GetAll(clientID uint) (promotion.Promotions, error) {
	s.setContext()

	promotions := promotion.Promotions{}
	err := s.db.Find(&promotions, "client_id = ?", clientID).Error
	if err != nil {
		return []promotion.Promotion{}, ErrNotFound
	}

	for i := 0; i < len(promotions); i++ {
		promotions[i].Days = promotion.SetDays(promotions[i].DaysString)
		promotions[i].DaysString = ""
		err := s.db.Where("type = ?", click.Promotion).
			Find(&promotions[i].Clicks, "type_id = ?", promotions[i].ID).Error
		if err != nil {
			return []promotion.Promotion{}, ErrNotFound
		}

		var ds DishStorage
		storedDish, err := ds.GetByID(promotions[i].DishID)
		if err != nil {
			continue
		}
		promotions[i].Dish = storedDish
	}

	return promotions, nil
}

// GetAllActive returns all stored promotions.
func (s PromotionStorage) GetAllActive(clientID uint) (promotion.Promotions, error) {
	s.setContext()

	promotions, err := s.GetAll(clientID)
	if err != nil {
		return []promotion.Promotion{}, err
	}

	result := promotion.Promotions{}
	now := time.Now()
	var cs ClientStorage
	c, err := cs.GetByID(clientID)
	if err != nil {
		return []promotion.Promotion{}, err
	}

	for _, p := range promotions {
		p.DaysString = promotion.SetDaysString(p.Days)
		if p.IsActive(now, c) {
			p.DaysString = ""
			result = append(result, p)
		}
	}

	return result, nil
}

// GetByID returns a promotion by ID.
func (s PromotionStorage) GetByID(id uint) (promotion.Promotion, error) {
	s.setContext()

	p := promotion.Promotion{}
	err := s.db.First(&p, "id = ?", id).Error
	if err != nil {
		return promotion.Promotion{}, ErrNotFound
	}

	p.Days = promotion.SetDays(p.DaysString)
	p.DaysString = ""

	err = s.db.Where("type = ?", click.Promotion).
		Find(&p.Clicks, "type_id = ?", p.ID).Error
	if err != nil {
		return promotion.Promotion{}, ErrNotFound
	}

	var ds DishStorage
	storedDish, _ := ds.GetByID(p.DishID)
	p.Dish = storedDish

	return p, nil
}
