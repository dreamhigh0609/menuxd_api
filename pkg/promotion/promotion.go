package promotion

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/super-radmir/menuxd_api/pkg/click"
	"github.com/super-radmir/menuxd_api/pkg/client"
	"github.com/super-radmir/menuxd_api/pkg/model"
)

// Storage handle Promotion's operations.
type Storage interface {
	Create(promotion *Promotion) error
	Update(id uint, promotion *Promotion) error
	Delete(id uint) error
	GetAll(clientID uint) (Promotions, error)
	GetAllActive(clientID uint) (Promotions, error)
	GetByID(id uint) (Promotion, error)
	AddClick(id uint) error
}

// Promotion is an client's event.
type Promotion struct {
	model.Model
	Title      string        `bson:"title" json:"title"`
	Pictures   []string      `gorm:"-" bson:"pictures" json:"pictures,omitempty"`
	PicturesString string    `gorm:"column:pictures" bson:"pictures" json:"pictures_string,omitempty"`
	StartAt    string        `bson:"start_at" json:"start_at"`
	EndAt      string        `bson:"end_at" json:"end_at"`
	Description string       `bson:"description,omitempty" json:"description,omitempty"`
	Price      float64       `bson:"price" json:"price"`
	Days       []string      `gorm:"-" bson:"days" json:"days"`
	Ingredients    []Ingredient `gorm:"-" json:"ingredients" json:"ingredients"`
	DaysString string        `gorm:"column:days" bson:"days" json:"days_string,omitempty"`
	Clicks     []click.Click `bson:"clicks" json:"clicks"`
	ClientID   uint          `bson:"client_id" json:"client_id"`
}

// Ingredient to the promotions.
type Ingredient struct {
	model.Model
	PromotionID  uint    `json:"promotion_id"`
	Price   float64 `json:"price"`
	Name    string  `json:"name"`
	OrderID *uint   `json:"order_id"`
	Active  bool    `json:"active"`
}

// SetString split strings into slices.
func SetString(arg []string) string {
	if len(arg) >= 0 {
		return strings.Join(arg, ",")
	}

	return ""
}

// SetSlice split strings into slices.
func SetSlice(arg string) []string {
	if arg != "" {
		return strings.Split(arg, ",")
	}

	return []string{"", "", ""}
}

func getDate(timeStr string, current time.Time, loc *time.Location) (t time.Time, err error) {
	timeArr := strings.Split(timeStr, ":")
	if len(timeArr) != 2 {
		return t, errors.New("Badly formatted time")
	}

	hStr, mStr := timeArr[0], timeArr[1]

	h, err := strconv.Atoi(hStr)
	if err != nil {
		return t, err
	}

	m, err := strconv.Atoi(mStr)
	if err != nil {
		return t, err
	}

	lt := time.Date(current.Year(), current.Month(), current.Day(), h, m, 0, 0, loc)

	return lt, nil
}

// SetDays split days string in slices.
func SetDays(s string) []string {
	if s != "" {
		return strings.Split(s, ",")
	}

	return []string{}
}

// SetDaysString join slice of days in a string.
func SetDaysString(s []string) string {
	if len(s) > 0 {
		return strings.Join(s, ",")
	}

	return ""
}

// IsActive confirm if the promotion is active.
func (p Promotion) IsActive(now time.Time, c client.Client) bool {
	loc, err := time.LoadLocation(c.Timezone)
	if err != nil {
		return false
	}

	currentTime := now.In(loc)

	currentDay := strings.ToLower(currentTime.Weekday().String())
	if ok := strings.Contains(p.DaysString, currentDay); !ok {
		return false
	}

	start, err := getDate(p.StartAt, currentTime, loc)
	if err != nil {
		return false
	}
	end, err := getDate(p.EndAt, currentTime, loc)
	if err != nil {
		return false
	}

	return now.After(start) && now.Before(end)
}

// Promotions alias for a slice of Promotions.
type Promotions []Promotion
