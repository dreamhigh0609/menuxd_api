package ad

import (
	"github.com/super-radmir/menuxd_api/pkg/click"
	"github.com/super-radmir/menuxd_api/pkg/model"
)

// Storage handle the CRUD operations with Ads
type Storage interface {
	Create(ad *Ad) error
	Update(id uint, ad *Ad) error
	Patch(id uint, updates map[string]interface{}) error
	Delete(id uint) error
	GetAll(clientID uint) (Ads, error)
	AddClick(adID uint) error
	GetByID(id uint) (Ad, error)
}

// Ad represents ads to the app.
type Ad struct {
	model.Model
	Title    string        `bson:"title" json:"title"`
	Active   bool          `bson:"active" json:"active"`
	Picture  string        `bson:"picture" json:"picture"`
	Clicks   []click.Click `bson:"clicks" json:"clicks"`
	ClientID uint          `bson:"client_id" json:"client_id"`
}

// Ads alias for the slice of Ads.
type Ads []Ad
