package controllertests

import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	"github.com/task/api/models"
)

func Database(databaseUrl string) error {

	var err error

	TestDbDriver := "postgres"

	if TestDbDriver == "postgres" {
		DBURL := databaseUrl
		server.DB, err = gorm.Open(TestDbDriver, DBURL)
		if err != nil {
			fmt.Printf("Cannot connect to %s database\n", TestDbDriver)
			return err
		} else {
			fmt.Printf("We are connected to the %s database\n", TestDbDriver)
		}
	}
	return err
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

	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}

	user := models.User{
		Username: "Pet",
		Email:    "pet@gmail.com",
		Password: "password",
	}

	err = server.DB.Model(&models.User{}).Create(&user).Error
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func seedUsers() ([]models.User, error) {

	var err error
	if err != nil {
		return nil, err
	}
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
			return []models.User{}, err
		}
	}
	return users, nil
}

func refreshUserAndPostTable() error {

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

func seedOneUserAndOnePost() (models.Product, error) {

	err := refreshUserAndPostTable()
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
	post := models.Product{
		ProductName:     "This is the title sam",
		AmountAvailable: 100,
		SellerID:        user.ID,
	}
	err = server.DB.Model(&models.Product{}).Create(&post).Error
	if err != nil {
		return models.Product{}, err
	}
	return post, nil
}

func seedUsersAndPosts() ([]models.User, []models.Product, error) {

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
	var posts = []models.Product{
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
		posts[i].SellerID = users[i].ID

		err = server.DB.Model(&models.Product{}).Create(&posts[i]).Error
		if err != nil {
			log.Fatalf("cannot seed posts table: %v", err)
		}
	}
	return users, posts, nil
}
