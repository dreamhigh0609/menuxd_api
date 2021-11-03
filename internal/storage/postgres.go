package storage

import (
	"os"
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/super-radmir/menuxd_api/pkg/ad"
	"github.com/super-radmir/menuxd_api/pkg/bill"
	"github.com/super-radmir/menuxd_api/pkg/category"
	"github.com/super-radmir/menuxd_api/pkg/click"
	"github.com/super-radmir/menuxd_api/pkg/client"
	"github.com/super-radmir/menuxd_api/pkg/dish"
	"github.com/super-radmir/menuxd_api/pkg/order"
	"github.com/super-radmir/menuxd_api/pkg/promotion"
	"github.com/super-radmir/menuxd_api/pkg/question"
	"github.com/super-radmir/menuxd_api/pkg/rating"
	"github.com/super-radmir/menuxd_api/pkg/stay"
	"github.com/super-radmir/menuxd_api/pkg/table"
	"github.com/super-radmir/menuxd_api/pkg/user"
	"github.com/super-radmir/menuxd_api/pkg/waiter"
)

// DBName Database name.
const DBName = "menuxd"

var (
	conn       *gorm.DB
	connString string
)

// createDBSession Create a new connection with the database.
func createDBSession() error {
	var err error
	conn, err = gorm.Open("postgres", connString)
	if err != nil {
		fmt.Println("test")
		fmt.Println(err)
		return err
	}
	fmt.Println("test2")
	fmt.Println(err)
	return nil
}

// getSession Returns the gorm conn.
func getSession() *gorm.DB {
	if conn == nil {
		createDBSession()
	}
	return conn
}

func migration() error {
	return conn.AutoMigrate(
		&ad.Ad{},
		&bill.Bill{},
		&category.Category{},
		&client.Client{},
		&dish.Dish{},
		&dish.Ingredient{},
		&promotion.Promotion{},
		&table.Table{},
		&user.User{},
		&waiter.Waiter{},
		&order.Order{},
		&order.Item{},
		&order.IngredientSelected{},
		&click.Click{},
		&stay.Stay{},
		&question.Question{},
		&rating.Rating{},
	).Error
}

// InitData initialize the conn.
func InitData() error {
	var err error
	err = createDBSession()
	if err != nil {
		return err
	}

	err = migration()
	if err != nil {
		return err
	}
	return nil
}

// Close the gorm conn.
func Close() error {
	return conn.Close()
}

func init() {
	connString = os.Getenv("DATABASE_URL")
}
