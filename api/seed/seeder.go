package seed

import (
	"log"

	"github.com/jinzhu/gorm"
	"github.com/task/api/models"
)

var users = []models.User{
	models.User{
		Username: "haitham rageh",
		Email:    "haitham@gmail.com",
		Password: "password",
		Role:     "seller",
		Deposit:  0,
	},
	models.User{
		Username: "osama zean",
		Email:    "osama@gmail.com",
		Password: "password",
		Role:     "buyer",
		Deposit:  1000,
	},
}

var Products = []models.Product{
	models.Product{
		ProductName:     "ProductName 1",
		AmountAvailable: 10,
		SellerID:        1,
		Price:           100,
	},
	models.Product{
		ProductName:     "ProductName 2",
		AmountAvailable: 20,
		SellerID:        1,
		Price:           200,
	},
}

func Load(db *gorm.DB) {

	err := db.Debug().DropTableIfExists(&models.Product{}, &models.User{}).Error
	if err != nil {
		log.Fatalf("cannot drop table: %v", err)
	}
	err = db.Debug().AutoMigrate(&models.User{}, &models.Product{}).Error
	if err != nil {
		log.Fatalf("cannot migrate table: %v", err)
	}

	err = db.Debug().Model(&models.Product{}).AddForeignKey("seller_id", "users(id)", "cascade", "cascade").Error
	if err != nil {
		log.Fatalf("attaching foreign key error: %v", err)
	}

	for i, _ := range users {
		err = db.Debug().Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			log.Fatalf("cannot seed users table: %v", err)
		}
		Products[i].SellerID = users[i].ID

		err = db.Debug().Model(&models.Product{}).Create(&Products[i]).Error
		if err != nil {
			log.Fatalf("cannot seed Products table: %v", err)
		}
	}
}
