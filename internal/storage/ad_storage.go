package storage

import (
	"github.com/jinzhu/gorm"
	"github.com/super-radmir/menuxd_api/pkg/ad"
	"github.com/super-radmir/menuxd_api/pkg/click"
)

// AdStorage storage to the ad model.
type AdStorage struct {
	session *Session
	db      *gorm.DB
}

// setContext initialize the context to AdStorage.
func (s *AdStorage) setContext() {
	s.session = NewSession()
	s.db = s.session.Client
}

// Create create a new ad.
func (s AdStorage) Create(a *ad.Ad) error {
	s.setContext()

	a.Active = true
	if a.Picture == "" {
		return ErrRequiredField
	}

	err := s.db.Create(a).Error
	if err != nil {
		return ErrNotInsert
	}

	return nil
}

// AddClick create a new click.
func (s AdStorage) AddClick(adID uint) error {
	s.setContext()

	c := click.Click{}
	c.TypeID = adID
	c.Type = click.Ad

	err := s.db.Create(&c).Error
	if err != nil {
		return ErrNotInsert
	}

	return nil
}

// Update update ad by ID.
func (s AdStorage) Update(id uint, a *ad.Ad) error {
	s.setContext()

	if a.Picture == "" {
		return ErrRequiredField
	}

	updates := map[string]interface{}{
		"picture": a.Picture,
		"title":   a.Title,
	}

	a.ID = id

	err := s.db.Model(a).Updates(updates).Error
	if err != nil {
		return ErrNotUpdate
	}

	return nil
}

// Patch update ad by ID.
func (s AdStorage) Patch(id uint, updates map[string]interface{}) error {
	s.setContext()

	delete(updates, "client_id")

	err := s.db.Model(&ad.Ad{}).Where("id = ?", id).Updates(updates).Error
	if err != nil {
		return ErrNotUpdate
	}

	return nil
}

// Delete remove an ad by ID.
func (s AdStorage) Delete(id uint) error {
	s.setContext()

	err := s.db.Delete(&ad.Ad{}, "id = ?", id).Error
	if err != nil {
		return ErrNotDelete
	}

	return nil
}

// GetAll returns all stored ads.
func (s AdStorage) GetAll(clientID uint) (ad.Ads, error) {
	s.setContext()

	ads := ad.Ads{}

	err := s.db.Find(&ads, "client_id = ?", clientID).Error
	if err != nil {
		return []ad.Ad{}, ErrNotFound
	}

	for i := 0; i < len(ads); i++ {
		err := s.db.Where("type = ?", click.Ad).
			Find(&ads[i].Clicks, "type_id = ?", ads[i].ID).Error
		if err != nil {
			return []ad.Ad{}, ErrNotFound
		}
	}

	return ads, nil
}

// GetByID returns an ad by ID.
func (s AdStorage) GetByID(id uint) (ad.Ad, error) {
	s.setContext()

	a := ad.Ad{}
	err := s.db.First(&a, "id = ?", id).Error
	if err != nil {
		return ad.Ad{}, ErrNotFound
	}

	return a, nil
}
