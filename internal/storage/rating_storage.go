package storage

import (
	"github.com/jinzhu/gorm"
	"github.com/super-radmir/menuxd_api/pkg/rating"
)

// RatingStorage storage to the rating model.
type RatingStorage struct {
	session *Session
	db      *gorm.DB
}

// setContext initialize the context to RatingStorage.
func (s *RatingStorage) setContext() {
	s.session = NewSession()
	s.db = s.session.Client
}

// Create create a new rating.
func (s RatingStorage) Create(r *rating.Rating) error {
	s.setContext()

	err := s.db.Create(r).Error
	if err != nil {
		return ErrNotInsert
	}

	return nil
}

// GetAll returns all stored ads.
func (s RatingStorage) GetAll() ([]rating.Rating, error) {
	s.setContext()

	ratings := []rating.Rating{}

	err := s.db.Find(&ratings).Error
	if err != nil {
		return []rating.Rating{}, ErrNotFound
	}

	return ratings, nil
}

// GetAllByQuestion returns all stored ads by question.
func (s RatingStorage) GetAllByQuestion(questionID uint) ([]rating.Rating, error) {
	s.setContext()

	ratings := []rating.Rating{}

	err := s.db.Find(&ratings, "question_id = ?", questionID).Error
	if err != nil {
		return []rating.Rating{}, ErrNotFound
	}

	return ratings, nil
}

// GetByID returns an rating by ID.
func (s RatingStorage) GetByID(id uint) (rating.Rating, error) {
	s.setContext()

	r := rating.Rating{}
	err := s.db.First(&r, "id = ?", id).Error
	if err != nil {
		return rating.Rating{}, ErrNotFound
	}

	return r, nil
}
