package controllertests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gorilla/mux"
	"github.com/task/api/models"
	"gopkg.in/go-playground/assert.v1"
)

func TestCreateProduct(t *testing.T) {

	err := refreshUserAndPostTable()
	if err != nil {
		log.Fatal(err)
	}
	user, err := seedOneUser()
	if err != nil {
		log.Fatalf("Cannot seed user %v\n", err)
	}
	token, err := server.SignIn(user.Email, "password") //Note the password in the database is already hashed, we want unhashed
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}
	tokenString := fmt.Sprintf("Bearer %v", token)

	samples := []struct {
		inputJSON       string
		statusCode      int
		proudctName     string
		amountAvailable int32
		seller_id       uint32
		tokenGiven      string
		errorMessage    string
	}{
		{
			inputJSON:       `{"proudct_name":"The proudct_name", "amount_available": 100, "seller_id": 1}`,
			statusCode:      201,
			tokenGiven:      tokenString,
			proudctName:     "The proudct_name",
			amountAvailable: 100,
			seller_id:       user.ID,
			errorMessage:    "",
		},
		{
			inputJSON:    `{"proudct_name":"The proudct_name", "amount_available": "the amount_available", "seller_id": 1}`,
			statusCode:   500,
			tokenGiven:   tokenString,
			errorMessage: "ProductName Already Taken",
		},
		{
			// When no token is passed
			inputJSON:    `{"proudct_name":"When no token is passed", "amount_available": "the amount_available", "seller_id": 1}`,
			statusCode:   401,
			tokenGiven:   "",
			errorMessage: "Unauthorized",
		},
		{
			// When incorrect token is passed
			inputJSON:    `{"proudct_name":"When incorrect token is passed", "amount_available": "the amount_available", "seller_id": 1}`,
			statusCode:   401,
			tokenGiven:   "This is an incorrect token",
			errorMessage: "Unauthorized",
		},
		{
			inputJSON:    `{"proudct_name": "", "amount_available": "The amount_available", "seller_id": 1}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required ProductName",
		},
		{
			inputJSON:    `{"proudct_name": "This is a proudct_name", "amount_available": "", "seller_id": 1}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required AmountAvailable",
		},
		{
			inputJSON:    `{"proudct_name": "This is an awesome proudct_name", "amount_available": "the amount_available"}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required Author",
		},
		{
			// When user 2 uses user 1 token
			inputJSON:    `{"proudct_name": "This is an awesome proudct_name", "amount_available": "the amount_available", "seller_id": 2}`,
			statusCode:   401,
			tokenGiven:   tokenString,
			errorMessage: "Unauthorized",
		},
	}
	for _, v := range samples {

		req, err := http.NewRequest("POST", "/posts", bytes.NewBufferString(v.inputJSON))
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.CreateProduct)

		req.Header.Set("Authorization", v.tokenGiven)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			fmt.Printf("Cannot convert to json: %v", err)
		}
		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 201 {
			assert.Equal(t, responseMap["proudct_name"], v.proudctName)
			assert.Equal(t, responseMap["amount_available"], v.amountAvailable)
			assert.Equal(t, responseMap["seller_id"], float64(v.seller_id)) //just for both ids to have the same type
		}
		if v.statusCode == 401 || v.statusCode == 422 || v.statusCode == 500 && v.errorMessage != "" {
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}

func TestGetPosts(t *testing.T) {

	err := refreshUserAndPostTable()
	if err != nil {
		log.Fatal(err)
	}
	_, _, err = seedUsersAndPosts()
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("GET", "/posts", nil)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.GetProducts)
	handler.ServeHTTP(rr, req)

	var posts []models.Product
	err = json.Unmarshal([]byte(rr.Body.String()), &posts)

	assert.Equal(t, rr.Code, http.StatusOK)
	assert.Equal(t, len(posts), 2)
}
func TestGetPostByID(t *testing.T) {

	err := refreshUserAndPostTable()
	if err != nil {
		log.Fatal(err)
	}
	post, err := seedOneUserAndOnePost()
	if err != nil {
		log.Fatal(err)
	}
	postSample := []struct {
		id              string
		statusCode      int
		proudct_name    string
		amountAvailable float32
		seller_id       uint32
		errorMessage    string
	}{
		{
			id:              strconv.Itoa(int(post.ID)),
			statusCode:      200,
			proudct_name:    post.ProductName,
			amountAvailable: post.AmountAvailable,
			seller_id:       post.SellerID,
		},
		{
			id:         "unknwon",
			statusCode: 400,
		},
	}
	for _, v := range postSample {

		req, err := http.NewRequest("GET", "/proudcts", nil)
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": v.id})

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.GetProduct)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			log.Fatalf("Cannot convert to json: %v", err)
		}
		assert.Equal(t, rr.Code, v.statusCode)

		if v.statusCode == 200 {
			assert.Equal(t, post.ProductName, responseMap["proudct_name"])
			assert.Equal(t, post.AmountAvailable, responseMap["amount_available"])
			assert.Equal(t, float64(post.SellerID), responseMap["seller_id"]) //the response author id is float64
		}
	}
}
