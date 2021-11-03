package notification

import (
	"time"

	"github.com/super-radmir/menuxd_api/pkg/table"
)

// Notification types.
const (
	CallWaiter uint = iota
	GetCheck
	MakeOrder
	Connected
)

// Notification is a message to send.
type Notification struct {
	Type     uint         `json:"type"`
	Message  string       `json:"message,omitempty"`
	Picture  string       `json:"picture"`
	Date     time.Time    `json:"date"`
	ClientID uint         `json:"clientId"`
	Active   bool         `json:"active"`
	Table    *table.Table `json:"table"`
}
