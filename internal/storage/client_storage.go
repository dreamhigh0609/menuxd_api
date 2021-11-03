package storage

import (
	"github.com/jinzhu/gorm"
	"github.com/super-radmir/menuxd_api/pkg/client"
	"github.com/super-radmir/menuxd_api/pkg/question"
)

// ClientStorage storage to the client model.
type ClientStorage struct {
	session *Session
	db      *gorm.DB
}

// setContext initialize the context to ClientStorage.
func (s *ClientStorage) setContext() {
	s.session = NewSession()
	s.db = s.session.Client
}

// Create create a new client.
func (s ClientStorage) Create(c *client.Client) error {
	s.setContext()

	if !c.ValidDate() {
		return ErrInvalidExpiration
	}

	if c.Name == "" {
		return ErrRequiredField
	}

	err := s.db.Create(c).Error
	if err != nil {
		return ErrNotInsert
	}

	var qs QuestionStorage

	q := question.Question{}
	q.ClientID = c.ID
	q.Text = "¿Qué le pareció la experiencia del Menu Digital?"
	q.Main = true

	err = qs.Create(&q)
	if err != nil {
		return ErrNotInsert
	}

	return nil
}

// Update update client by ID.
func (s ClientStorage) Update(id uint, c *client.Client) error {
	s.setContext()

	if !c.ValidDate() {
		return ErrInvalidExpiration
	}

	updates := map[string]interface{}{
		"name":      c.Name,
		"picture":   c.Picture,
		"active":    c.Active,
		"expire_at": c.ExpireAt,
	}

	err := s.db.Model(&client.Client{}).Where("id = ?", id).Updates(updates).Error
	if err != nil {
		return ErrNotUpdate
	}

	return nil
}

// Delete remove a client by ID.
func (s ClientStorage) Delete(id uint) error {
	s.setContext()

	err := s.db.Delete(&client.Client{}, "id = ?", id).Error
	if err != nil {
		return ErrNotDelete
	}

	return nil
}

// GetAll returns all stored clients.
func (s ClientStorage) GetAll(userID uint) (client.Clients, error) {
	s.setContext()

	clients := client.Clients{}

	err := s.db.Find(&clients, "user_id = ?", userID).Error
	if err != nil {
		return []client.Client{}, ErrNotFound
	}

	return clients, nil
}

// GetByID returns a client by ID.
func (s ClientStorage) GetByID(id uint) (client.Client, error) {
	s.setContext()

	c := client.Client{}
	err := s.db.First(&c, "id = ?", id).Error
	if err != nil {
		return client.Client{}, ErrNotFound
	}

	return c, nil
}
