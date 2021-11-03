package storage

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/sethvargo/go-password/password"
	"github.com/super-radmir/menuxd_api/pkg/email"
	"github.com/super-radmir/menuxd_api/pkg/user"
)

// UserStorage storage to the user model
type UserStorage struct {
	session *Session
	db      *gorm.DB
}

// setContext initialize the context to UserStorage
func (s *UserStorage) setContext() {
	s.session = NewSession()
	s.db = s.session.Client
}

// Create create a new user
func (s UserStorage) Create(u *user.User) error {
	s.setContext()

	pass, err := password.Generate(10, 3, 0, false, false)
	if err != nil {
		fmt.Println(err)
		return ErrPasswordGeneration
	}

	u.Password = pass
	u.PreparePass()

	if u.ImageURL == "" {
		u.SetGravatar()
	}

	if ok := u.ValidateEmail(); !ok {
		fmt.Println("Email invalid")
		return ErrNotInsert
	}

	err = s.db.Create(u).Error
	if err != nil {
		fmt.Println(err)
		return ErrNotInsert
	}

	e := email.NewUser(u.ID, u.Email, pass)
	if err = e.Send(); err != nil {
		fmt.Println(err)
		return ErrSendEmail
	}
	return nil
}

// Update update user by ID
func (s UserStorage) Update(id uint, u *user.User) error {
	s.setContext()
	u.ID = id

	updates := map[string]interface{}{
		"role":      u.Role,
		"image_url": u.ImageURL,
	}

	err := s.db.Model(u).Updates(updates).Error
	if err != nil {
		return ErrNotUpdate
	}

	return nil
}

// Delete remove a user by ID
func (s UserStorage) Delete(id uint) error {
	s.setContext()

	err := s.db.Delete(&user.User{}, "id = ?", id).Error
	if err != nil {
		return ErrNotDelete
	}

	return nil
}

// GetAll returns all stored users
func (s UserStorage) GetAll() (user.Users, error) {
	s.setContext()

	users := user.Users{}
	err := s.db.Find(&users).Error
	if err != nil {
		return []user.User{}, ErrNotFound
	}

	for i := 0; i < len(users); i++ {
		users[i].CleanPass()
	}

	return users, nil
}

// GetByID returns a user by ID
func (s UserStorage) GetByID(id uint) (user.User, error) {
	s.setContext()

	u := user.User{}
	err := s.db.First(&u, id).Error
	if err != nil {
		return user.User{}, ErrNotFound
	}
	u.CleanPass()

	return u, nil
}

// GetByEmail returns a user by email address
func (s UserStorage) GetByEmail(email string) (user.User, error) {
	s.setContext()

	u := user.User{}
	err := s.db.First(&u, "email = ?", email).Error
	if err != nil {
		return user.User{}, ErrNotFound
	}

	return u, nil
}

// Confirm change user state confirmed to true and set the new password
func (s UserStorage) Confirm(id uint, u *user.User) error {
	s.setContext()

	if confirmed := u.ConfirmPass(); !confirmed {
		return ErrNotMatch
	}

	err := u.PreparePass()
	if err != nil {
		return err
	}

	u.ID = id
	err = s.db.Model(u).Updates(map[string]interface{}{
		"confirmed":    true,
		"HashPassword": u.HashPassword,
	}).Error

	if err != nil {
		return ErrNotUpdate
	}

	return nil
}

// RecoverPassword generate a temporal password and send a email to the user
// with the new password
func (s UserStorage) RecoverPassword(address string) error {
	s.setContext()

	pass, err := password.Generate(10, 3, 0, false, false)
	if err != nil {
		return ErrPasswordGeneration
	}

	u := user.New()

	u.Password = pass
	u.PreparePass()

	err = s.db.Model(&user.User{}).Where("email = ?", address).Update("HashPassword", u.HashPassword).Error
	if err != nil {
		return ErrNotUpdate
	}

	if err != nil {
		return ErrNotFound
	}

	e := email.ChangePassword(address, pass)
	if err = e.Send(); err != nil {
		return ErrSendEmail
	}

	return nil
}
