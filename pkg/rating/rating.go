package rating

import (
	"errors"

	"github.com/super-radmir/menuxd_api/pkg/model"
)

// Storage handle the CRUD operations with Rating.
type Storage interface {
	Create(r *Rating) error
	GetAll() ([]Rating, error)
	GetByID(id uint) (Rating, error)
	GetAllByQuestion(questionID uint) ([]Rating, error)
}

// Rating is a customer score for services.
type Rating struct {
	model.Model
	QuestionID uint `bson:"question_id" json:"question_id,omitempty"`
	Score      uint `bson:"score" json:"score"`
}

// IsValid verifies that the scores are valid.
func (r Rating) IsValid() bool {
	return r.Score > 0 && r.Score <= 5
}

// BeforeSave checks that the category is valid before inserting.
func (r *Rating) BeforeSave() (err error) {
	if !r.IsValid() {
		err = errors.New("can't save invalid data")
	}

	return
}
