package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/task/api/auth"
	"github.com/task/api/models"
	"github.com/task/api/responses"
	"github.com/task/api/utils/formaterror"
	"golang.org/x/crypto/bcrypt"
)

func (server *Server) Login(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	type Userlogin struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	userlogin := Userlogin{}
	err = json.Unmarshal(body, &userlogin)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	user := models.User{}
	user.Email = userlogin.Email
	user.Password = userlogin.Password
	user.Prepare()
	err = user.Validate("login")
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	token, err := server.SignIn(user.Email, user.Password)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusUnprocessableEntity, formattedError)
		return
	}
	type Loginresp struct {
		Jwt string `json:"jwt"`
	}
	login := new(Loginresp)
	login.Jwt = token
	responses.JSON(w, http.StatusOK, login)
}

func (server *Server) SignIn(email, password string) (string, error) {

	var err error

	user := models.User{}

	err = server.DB.Debug().Model(models.User{}).Where("email = ?", email).Take(&user).Error
	if err != nil {
		return "", err
	}
	err = models.VerifyPassword(user.Password, password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", err
	}
	return auth.CreateToken(user.ID, user.Role)
}
