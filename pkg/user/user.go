package user

import (
	"crypto/md5"
	"fmt"
	"regexp"

	"github.com/super-radmir/menuxd_api/pkg/model"

	"golang.org/x/crypto/bcrypt"
)

// Storage handle the CRUD operations with Users.
type Storage interface {
	Create(user *User) error
	Update(id uint, user *User) error
	Confirm(id uint, user *User) error
	Delete(id uint) error
	GetAll() (Users, error)
	GetByID(id uint) (User, error)
	GetByEmail(email string) (User, error)
	RecoverPassword(address string) error
}

// User of the system.
type User struct {
	model.Model
	Email           string `gorm:"unique_index" json:"email"`
	Password        string `gorm:"-" bson:"-" json:"password,omitempty"`
	Role            string `gorm:"default:'client'" json:"role,omitempty"`
	ImageURL        string `bson:"image_url" json:"image_url,omitempty"`
	Confirmed       bool   `gorm:"default:false" json:"confirmed,omitempty"`
	ConfirmPassword string `gorm:"-" bson:"-" json:"confirm_password,omitempty"`
	OldPassword     string `gorm:"-" bson:"-" json:"old_password,omitempty"`
	HashPassword    string `gorm:"column:password" bson:"password" json:"hash_password,omitempty"`
}

// ConfirmPass compare Password and ConfirmPassword and return true if they are the same
func (u User) ConfirmPass() bool {
	return u.Password == u.ConfirmPassword
}

// ComparePass compare de HashPassword with a raw password and return true if they are the same
func (u User) ComparePass(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.HashPassword), []byte(password))
	if err != nil {
		return false
	}
	return true
}

// PreparePass generate a hashed with the password and put the result in HashPassword
func (u *User) PreparePass() error {
	hpass, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.HashPassword = string(hpass)
	return nil
}

func (u User) ValidateEmail() bool {
	const pattern = "^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"

	re := regexp.MustCompile(pattern)

	return re.MatchString(u.Email)
}

// SetGravatar generate a url for Gravatar and set the result in ImageURL
func (u *User) SetGravatar() {
	picmd5 := md5.Sum([]byte(u.Email))
	u.ImageURL = fmt.Sprintf("https://gravatar.com/avatar/%x?s=100", picmd5)
}

// CleanPass reset the values of Password and ConfirmPassword to empty string
func (u *User) CleanPass() {
	u.Password = ""
	u.ConfirmPassword = ""
	u.HashPassword = ""
	u.OldPassword = ""
}

// Users alias for a slice of Users
type Users []User

// New returns a instance of User with default configuration
func New() *User {
	return &User{}
}
