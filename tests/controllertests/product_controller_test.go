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
		title           string
		amountAvailable int32
		seller_id       uint32
		tokenGiven      string
		errorMessage    string
	}{
		{
			inputJSON:       `{"title":"The title", "amount_available": 100, "seller_id": 1}`,
			statusCode:      201,
			tokenGiven:      tokenString,
			title:           "The title",
			amountAvailable: 100,
			seller_id:       user.ID,
			errorMessage:    "",
		},
		{
			inputJSON:    `{"title":"The title", "content": "the content", "seller_id": 1}`,
			statusCode:   500,
			tokenGiven:   tokenString,
			errorMessage: "ProductName Already Taken",
		},
		{
			// When no token is passed
			inputJSON:    `{"title":"When no token is passed", "content": "the content", "seller_id": 1}`,
			statusCode:   401,
			tokenGiven:   "",
			errorMessage: "Unauthorized",
		},
		{
			// When incorrect token is passed
			inputJSON:    `{"title":"When incorrect token is passed", "content": "the content", "seller_id": 1}`,
			statusCode:   401,
			tokenGiven:   "This is an incorrect token",
			errorMessage: "Unauthorized",
		},
		{
			inputJSON:    `{"title": "", "content": "The content", "seller_id": 1}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required ProductName",
		},
		{
			inputJSON:    `{"title": "This is a title", "content": "", "seller_id": 1}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required AmountAvailable",
		},
		{
			inputJSON:    `{"title": "This is an awesome title", "content": "the content"}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required Author",
		},
		{
			// When user 2 uses user 1 token
			inputJSON:    `{"title": "This is an awesome title", "content": "the content", "seller_id": 2}`,
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
			assert.Equal(t, responseMap["title"], v.title)
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
		title           string
		amountAvailable float32
		seller_id       uint32
		errorMessage    string
	}{
		{
			id:              strconv.Itoa(int(post.ID)),
			statusCode:      200,
			title:           post.ProductName,
			amountAvailable: post.AmountAvailable,
			seller_id:       post.SellerID,
		},
		{
			id:         "unknwon",
			statusCode: 400,
		},
	}
	for _, v := range postSample {

		req, err := http.NewRequest("GET", "/posts", nil)
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
			assert.Equal(t, post.ProductName, responseMap["title"])
			assert.Equal(t, post.AmountAvailable, responseMap["amount_available"])
			assert.Equal(t, float64(post.SellerID), responseMap["seller_id"]) //the response author id is float64
		}
	}
}

func TestUpdatePost(t *testing.T) {

	var PostUserEmail, PostUserPassword string
	var AuthPostSellerID uint32
	var AuthPostID uint64

	err := refreshUserAndPostTable()
	if err != nil {
		log.Fatal(err)
	}
	users, posts, err := seedUsersAndPosts()
	if err != nil {
		log.Fatal(err)
	}
	// Get only the first user
	for _, user := range users {
		if user.ID == 2 {
			continue
		}
		PostUserEmail = user.Email
		PostUserPassword = "password" //Note the password in the database is already hashed, we want unhashed
	}
	//Login the user and get the authentication token
	token, err := server.SignIn(PostUserEmail, PostUserPassword)
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}
	tokenString := fmt.Sprintf("Bearer %v", token)

	// Get only the first post
	for _, post := range posts {
		if post.ID == 2 {
			continue
		}
		AuthPostID = post.ID
		AuthPostSellerID = post.SellerID
	}
	// fmt.Printf("this is the auth post: %v\n", AuthPostID)

	samples := []struct {
		id              string
		updateJSON      string
		statusCode      int
		title           string
		amountAvailable float32
		seller_id       uint32
		tokenGiven      string
		errorMessage    string
	}{
		{
			// Convert int64 to int first before converting to string
			id:              strconv.Itoa(int(AuthPostID)),
			updateJSON:      `{"title":"The updated post", "amount_available": "100", "seller_id": 1}`,
			statusCode:      200,
			title:           "The updated post",
			amountAvailable: 100,
			seller_id:       AuthPostSellerID,
			tokenGiven:      tokenString,
			errorMessage:    "",
		},
		{
			// When no token is provided
			id:           strconv.Itoa(int(AuthPostID)),
			updateJSON:   `{"title":"This is still another title", "content": "This is the updated content", "seller_id": 1}`,
			tokenGiven:   "",
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
		{
			// When incorrect token is provided
			id:           strconv.Itoa(int(AuthPostID)),
			updateJSON:   `{"title":"This is still another title", "content": "This is the updated content", "seller_id": 1}`,
			tokenGiven:   "this is an incorrect token",
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
		{
			//Note: "ProductName 2" belongs to post 2, and title must be unique
			id:           strconv.Itoa(int(AuthPostID)),
			updateJSON:   `{"title":"ProductName 2", "content": "This is the updated content", "seller_id": 1}`,
			statusCode:   500,
			tokenGiven:   tokenString,
			errorMessage: "ProductName Already Taken",
		},
		{
			id:           strconv.Itoa(int(AuthPostID)),
			updateJSON:   `{"title":"", "content": "This is the updated content", "seller_id": 1}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required ProductName",
		},
		{
			id:           strconv.Itoa(int(AuthPostID)),
			updateJSON:   `{"title":"Awesome title", "content": "", "seller_id": 1}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required AmountAvailable",
		},
		{
			id:           strconv.Itoa(int(AuthPostID)),
			updateJSON:   `{"title":"This is another title", "content": "This is the updated content"}`,
			statusCode:   401,
			tokenGiven:   tokenString,
			errorMessage: "Unauthorized",
		},
		{
			id:         "unknwon",
			statusCode: 400,
		},
		{
			id:           strconv.Itoa(int(AuthPostID)),
			updateJSON:   `{"title":"This is still another title", "content": "This is the updated content", "seller_id": 2}`,
			tokenGiven:   tokenString,
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
	}

	for _, v := range samples {

		req, err := http.NewRequest("POST", "/posts", bytes.NewBufferString(v.updateJSON))
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": v.id})
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.UpdateProduct)

		req.Header.Set("Authorization", v.tokenGiven)

		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			t.Errorf("Cannot convert to json: %v", err)
		}
		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 200 {
			assert.Equal(t, responseMap["title"], v.title)
			assert.Equal(t, responseMap["amount_available"], v.amountAvailable)
			assert.Equal(t, responseMap["seller_id"], float64(v.seller_id)) //just to match the type of the json we receive thats why we used float64
		}
		if v.statusCode == 401 || v.statusCode == 422 || v.statusCode == 500 && v.errorMessage != "" {
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}

func TestDeletePost(t *testing.T) {

	var PostUserEmail, PostUserPassword string
	var PostUserID uint32
	var AuthPostID uint64

	err := refreshUserAndPostTable()
	if err != nil {
		log.Fatal(err)
	}
	users, posts, err := seedUsersAndPosts()
	if err != nil {
		log.Fatal(err)
	}
	//Let's get only the Second user
	for _, user := range users {
		if user.ID == 1 {
			continue
		}
		PostUserEmail = user.Email
		PostUserPassword = "password" //Note the password in the database is already hashed, we want unhashed
	}
	//Login the user and get the authentication token
	token, err := server.SignIn(PostUserEmail, PostUserPassword)
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}
	tokenString := fmt.Sprintf("Bearer %v", token)

	// Get only the second post
	for _, post := range posts {
		if post.ID == 1 {
			continue
		}
		AuthPostID = post.ID
		PostUserID = post.SellerID
	}
	postSample := []struct {
		id           string
		seller_id    uint32
		tokenGiven   string
		statusCode   int
		errorMessage string
	}{
		{
			// Convert int64 to int first before converting to string
			id:           strconv.Itoa(int(AuthPostID)),
			seller_id:    PostUserID,
			tokenGiven:   tokenString,
			statusCode:   204,
			errorMessage: "",
		},
		{
			// When empty token is passed
			id:           strconv.Itoa(int(AuthPostID)),
			seller_id:    PostUserID,
			tokenGiven:   "",
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
		{
			// When incorrect token is passed
			id:           strconv.Itoa(int(AuthPostID)),
			seller_id:    PostUserID,
			tokenGiven:   "This is an incorrect token",
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
		{
			id:         "unknwon",
			tokenGiven: tokenString,
			statusCode: 400,
		},
		{
			id:           strconv.Itoa(int(1)),
			seller_id:    1,
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
	}
	for _, v := range postSample {

		req, _ := http.NewRequest("GET", "/posts", nil)
		req = mux.SetURLVars(req, map[string]string{"id": v.id})

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.DeleteProduct)

		req.Header.Set("Authorization", v.tokenGiven)

		handler.ServeHTTP(rr, req)

		assert.Equal(t, rr.Code, v.statusCode)

		if v.statusCode == 401 && v.errorMessage != "" {

			responseMap := make(map[string]interface{})
			err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
			if err != nil {
				t.Errorf("Cannot convert to json: %v", err)
			}
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}
