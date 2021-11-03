package client

import (
	"time"

	"github.com/super-radmir/menuxd_api/pkg/model"
)

// Storage handle the CRUD operations with Clients.
type Storage interface {
	Create(client *Client) error
	Update(id uint, client *Client) error
	Delete(id uint) error
	GetAll(userID uint) (Clients, error)
	GetByID(id uint) (Client, error)
}

// Client is a restaurant to MenuXD system.
type Client struct {
	model.Model
	Name     string    `json:"name"`
	Picture  string    `gorm:"default:'http://localhost:8080/public/client-default.png'" json:"picture,omitempty"`
	Active   bool      `gorm:"default:true" json:"active,omitempty"`
	UserID   uint      `bson:"user_id" json:"user_id"`
	Timezone string    `gorm:"default:'America/Asuncion'" bson:"timezone" json:"timezone"`
	ExpireAt time.Time `bson:"expire_at" json:"expire_at,omitempty"`
}

// ValidDate confirm the date to expire the client.
func (c Client) ValidDate() bool {
	return c.ExpireAt.After(time.Now())
}

// Clients alias for a slice of Clients.
type Clients []Client

// New returns a instance of Client with default configuration.
func New() Client {
	return Client{}
}
