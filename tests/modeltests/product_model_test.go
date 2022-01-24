package modeltests

import (
	"log"
	"testing"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/task/api/models"
	"gopkg.in/go-playground/assert.v1"
)

func TestFindAllProducts(t *testing.T) {

	err := refreshUserAndProductTable()
	if err != nil {
		log.Fatalf("Error refreshing user and post table %v\n", err)
	}
	_, _, err = seedUsersAndProducts()
	if err != nil {
		log.Fatalf("Error seeding user and post  table %v\n", err)
	}
	products, err := productInstance.FindAllProducts(server.DB)
	if err != nil {
		t.Errorf("this is the error getting the products: %v\n", err)
		return
	}
	assert.Equal(t, len(*products), 2)
}

func TestSaveProduct(t *testing.T) {

	err := refreshUserAndProductTable()
	if err != nil {
		log.Fatalf("Error user and post refreshing table %v\n", err)
	}

	user, err := seedOneUser()
	if err != nil {
		log.Fatalf("Cannot seed user %v\n", err)
	}

	newProduct := models.Product{
		ID:              1,
		ProductName:     "This is the title",
		AmountAvailable: 100,
		SellerID:        user.ID,
	}
	savedProduct, err := newProduct.SaveProduct(server.DB)
	if err != nil {
		t.Errorf("this is the error getting the post: %v\n", err)
		return
	}
	assert.Equal(t, newProduct.ID, savedProduct.ID)
	assert.Equal(t, newProduct.ProductName, savedProduct.ProductName)
	assert.Equal(t, newProduct.AmountAvailable, savedProduct.AmountAvailable)
	assert.Equal(t, newProduct.SellerID, savedProduct.SellerID)

}

func TestGetProductByID(t *testing.T) {

	err := refreshUserAndProductTable()
	if err != nil {
		log.Fatalf("Error refreshing user and post table: %v\n", err)
	}
	post, err := seedOneUserAndOneProduct()
	if err != nil {
		log.Fatalf("Error Seeding table")
	}
	foundProduct, err := productInstance.FindProductByID(server.DB, post.ID)
	if err != nil {
		t.Errorf("this is the error getting one user: %v\n", err)
		return
	}
	assert.Equal(t, foundProduct.ID, post.ID)
	assert.Equal(t, foundProduct.ProductName, post.ProductName)
	assert.Equal(t, foundProduct.AmountAvailable, post.AmountAvailable)
}

func TestUpdateAProduct(t *testing.T) {

	err := refreshUserAndProductTable()
	if err != nil {
		log.Fatalf("Error refreshing user and post table: %v\n", err)
	}
	post, err := seedOneUserAndOneProduct()
	if err != nil {
		log.Fatalf("Error Seeding table")
	}
	proudctUpdate := models.Product{
		ID:              1,
		ProductName:     "modiUpdate",
		AmountAvailable: 100,
		SellerID:        post.SellerID,
	}
	updatedProduct, err := proudctUpdate.UpdateAProduct(server.DB)
	if err != nil {
		t.Errorf("this is the error updating the user: %v\n", err)
		return
	}
	assert.Equal(t, updatedProduct.ID, proudctUpdate.ID)
	assert.Equal(t, updatedProduct.ProductName, proudctUpdate.ProductName)
	assert.Equal(t, updatedProduct.AmountAvailable, proudctUpdate.AmountAvailable)
	assert.Equal(t, updatedProduct.SellerID, proudctUpdate.SellerID)
}

func TestDeleteAProduct(t *testing.T) {

	err := refreshUserAndProductTable()
	if err != nil {
		log.Fatalf("Error refreshing user and post table: %v\n", err)
	}
	post, err := seedOneUserAndOneProduct()
	if err != nil {
		log.Fatalf("Error Seeding tables")
	}
	isDeleted, err := productInstance.DeleteAProduct(server.DB, post.ID, post.SellerID)
	if err != nil {
		t.Errorf("this is the error updating the user: %v\n", err)
		return
	}
	//one shows that the record has been deleted or:
	// assert.Equal(t, int(isDeleted), 1)

	//Can be done this way too
	assert.Equal(t, isDeleted, int64(1))
}
