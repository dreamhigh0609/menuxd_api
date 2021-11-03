package question

import (
	"github.com/super-radmir/menuxd_api/pkg/model"
	"github.com/super-radmir/menuxd_api/pkg/rating"
)

// Storage handle the CRUD operations with Questions.
type Storage interface {
	Create(question *Question) error
	Update(id uint, question *Question) error
	Delete(id uint) error
	GetAll(clientID uint) ([]Question, error)
	GetByID(id uint) (Question, error)
}

// Question is a client question to ask customers.
type Question struct {
	model.Model
	Text     string          `bson:"text" json:"text"`
	Main     bool            `bson:"main" json:"main"`
	Rating   []rating.Rating `bson:"rating" json:"rating"`
	ClientID uint            `bson:"client_id" json:"client_id,omitempty"`
}
