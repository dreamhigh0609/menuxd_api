package storage

import (
	"github.com/jinzhu/gorm"

	"github.com/super-radmir/menuxd_api/pkg/question"
)

// QuestionStorage storage to the question model.
type QuestionStorage struct {
	session *Session
	db      *gorm.DB
}

// setContext initialize the context to QuestionStorage.
func (s *QuestionStorage) setContext() {
	s.session = NewSession()
	s.db = s.session.Client
}

// Create create a new question.
func (s QuestionStorage) Create(q *question.Question) error {
	s.setContext()

	if q.Text == "" {
		return ErrNotInsert
	}

	q.Main = false

	err := s.db.Create(q).Error
	if err != nil {
		return ErrNotInsert
	}
	return nil
}

// Update update a question by ID.
func (s QuestionStorage) Update(id uint, q *question.Question) error {
	s.setContext()

	q.ID = id

	if q.Main {
		return ErrNotInsert
	}

	err := s.db.Model(q).Update("text", q.Text).Error
	if err != nil {
		return ErrNotUpdate
	}

	return nil
}

// Delete remove a question by ID.
func (s QuestionStorage) Delete(id uint) error {
	s.setContext()

	err := s.db.Delete(&question.Question{}, "id = ?", id).Error
	if err != nil {
		return ErrNotDelete
	}

	return nil
}

// GetAll returns all stored questions.
func (s QuestionStorage) GetAll(clientID uint) ([]question.Question, error) {
	s.setContext()

	questions := []question.Question{}
	err := s.db.Order("text").Find(&questions, "client_id = ?", clientID).Error
	if err != nil {
		return []question.Question{}, ErrNotFound
	}

	var rs RatingStorage
	result := []question.Question{}
	for _, q := range questions {
		ratings, _ := rs.GetAllByQuestion(q.ID)
		q.Rating = ratings
		result = append(result, q)
	}

	return result, nil
}

// GetByID returns a question by ID.
func (s QuestionStorage) GetByID(id uint) (question.Question, error) {
	s.setContext()

	q := question.Question{}
	err := s.db.First(&q, "id = ?", id).Error
	if err != nil {
		return question.Question{}, ErrNotFound
	}
	var rs RatingStorage
	ratings, _ := rs.GetAllByQuestion(q.ID)
	q.Rating = ratings

	return q, nil
}
