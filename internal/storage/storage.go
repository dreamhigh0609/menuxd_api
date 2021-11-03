package storage

import (
	"errors"

	"github.com/jinzhu/gorm"
)

// Errors.
var (
	ErrPasswordGeneration = errors.New("failed password generation")
	ErrNotInsert          = errors.New("could not insert it")
	ErrNotUpdate          = errors.New("could not update it")
	ErrNotDelete          = errors.New("could not delete it")
	ErrNotFound           = errors.New("could not found it")
	ErrParseID            = errors.New("could not parse id")
	ErrSendEmail          = errors.New("failed to send email")
	ErrNotMatch           = errors.New("passwords don't match")
	ErrInvalidExpiration  = errors.New("invalid expiration date")
	ErrBadRequest         = errors.New("bad request")
	ErrInvalidHourRange   = errors.New("invalid hour range")
	ErrInvalidPIN         = errors.New("pin invalid")
	ErrRequiredField      = errors.New("required field")
)

// Session is the Storage session.
type Session struct {
	Client *gorm.DB
}

// Close the client.
func (s *Session) Close() error {
	return s.Client.Close()
}

// NewSession returns a new session.
func NewSession() *Session {
	return &Session{
		Client: getSession(),
	}
}
