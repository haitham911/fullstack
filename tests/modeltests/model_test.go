package modeltests

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"github.com/task/api/controllers"
	"github.com/task/api/models"
)

var server = controllers.Server{}
var userInstance = models.User{}
var productInstance = models.Product{}

func TestMain(m *testing.M) {
	var err error
	err = godotenv.Load(os.ExpandEnv("../../.env"))
	if err != nil {
		log.Fatalf("Error getting env %v\n", err)
	}
	MainDB()
	time.Sleep(30 * time.Second)

	os.Exit(m.Run())
}

func Database(databaseUrl string) {

	var err error

	TestDbDriver := os.Getenv("TestDbDriver")

	if TestDbDriver == "postgres" {
		DBURL := databaseUrl
		server.DB, err = gorm.Open(TestDbDriver, DBURL)
		if err != nil {
			fmt.Printf("Cannot connect to %s database\n", TestDbDriver)
			log.Fatal("This is the error:", err)
		} else {
			fmt.Printf("We are connected to the %s database\n", TestDbDriver)
		}
	}
}

func refreshUserTable() error {
	err := server.DB.DropTableIfExists(&models.User{}).Error
	if err != nil {
		return err
	}
	err = server.DB.AutoMigrate(&models.User{}).Error
	if err != nil {
		return err
	}
	log.Printf("Successfully refreshed table")
	return nil
}

func seedOneUser() (models.User, error) {

	refreshUserTable()

	user := models.User{
		Username: "Pet",
		Email:    "pet@gmail.com",
		Password: "password",
	}

	err := server.DB.Model(&models.User{}).Create(&user).Error
	if err != nil {
		log.Fatalf("cannot seed users table: %v", err)
	}
	return user, nil
}

func seedUsers() error {

	users := []models.User{
		models.User{
			Username: "Steven victor",
			Email:    "steven@gmail.com",
			Password: "password",
		},
		models.User{
			Username: "Kenny Morris",
			Email:    "kenny@gmail.com",
			Password: "password",
		},
	}

	for i, _ := range users {
		err := server.DB.Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func refreshUserAndProductTable() error {

	err := server.DB.DropTableIfExists(&models.User{}, &models.Product{}).Error
	if err != nil {
		return err
	}
	err = server.DB.AutoMigrate(&models.User{}, &models.Product{}).Error
	if err != nil {
		return err
	}
	log.Printf("Successfully refreshed tables")
	return nil
}

func seedOneUserAndOneProduct() (models.Product, error) {

	err := refreshUserAndProductTable()
	if err != nil {
		return models.Product{}, err
	}
	user := models.User{
		Username: "Sam Phil",
		Email:    "sam@gmail.com",
		Password: "password",
	}
	err = server.DB.Model(&models.User{}).Create(&user).Error
	if err != nil {
		return models.Product{}, err
	}
	product := models.Product{
		ProductName:     "This is the title sam",
		AmountAvailable: 1000,
		SellerID:        user.ID,
	}
	err = server.DB.Model(&models.Product{}).Create(&product).Error
	if err != nil {
		return models.Product{}, err
	}
	return product, nil
}

func seedUsersAndProducts() ([]models.User, []models.Product, error) {

	var err error

	if err != nil {
		return []models.User{}, []models.Product{}, err
	}
	var users = []models.User{
		models.User{
			Username: "Steven victor",
			Email:    "steven@gmail.com",
			Password: "password",
		},
		models.User{
			Username: "Magu Frank",
			Email:    "magu@gmail.com",
			Password: "password",
		},
	}
	var products = []models.Product{
		models.Product{
			ProductName:     "ProductName 1",
			AmountAvailable: 100,
		},
		models.Product{
			ProductName:     "ProductName 2",
			AmountAvailable: 100,
		},
	}

	for i, _ := range users {
		err = server.DB.Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			log.Fatalf("cannot seed users table: %v", err)
		}
		products[i].SellerID = users[i].ID

		err = server.DB.Model(&models.Product{}).Create(&products[i]).Error
		if err != nil {
			log.Fatalf("cannot seed Products table: %v", err)
		}
	}
	return users, products, nil
}
