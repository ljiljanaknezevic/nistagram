package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"user-service-mod/model"
	"user-service-mod/service"

	"github.com/gorilla/mux"
)

type UserHandler struct {
	Service *service.UserService
}

type EmailForRecovery struct {
	Email string `json:"email"`
}

type ChangePassword struct {
	Email        string `json:"email"`
	Password     string `json:"newPass"`
	ConfPassword string `json:"confirmPass"`
}

func (handler *UserHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	b, _ := ioutil.ReadAll(r.Body)
	var changePassword ChangePassword
	err := json.Unmarshal(b, &changePassword)
	if err != nil {
		var err model.Error
		err = model.SetError(err, "Error in reading payload.")
		json.NewEncoder(w).Encode(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var user model.User
	user = handler.Service.GetUserByEmailAddress(changePassword.Email)
	if (model.User{}) == user {
		var err model.Error
		err = model.SetError(err, "User with that email doesn't exist.")

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}
	user.Password, _ = handler.Service.GeneratehashPassword(changePassword.Password)
	handler.Service.UpdateUser(&user)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
}

func (handler *UserHandler) ChangeUserData(w http.ResponseWriter, r *http.Request) {
	b, _ := ioutil.ReadAll(r.Body)
	var userChange model.User
	err := json.Unmarshal(b, &userChange)
	if err != nil {
		var err model.Error
		err = model.SetError(err, "Error in reading payload.")
		json.NewEncoder(w).Encode(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var user model.User
	user = handler.Service.GetUserByEmailAddress(userChange.Email)
	if (model.User{}) == user {
		var err model.Error
		err = model.SetError(err, "User with that email doesn't exist.")

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}
	user.Email = userChange.Email
	user.Username = userChange.Username
	user.Name = userChange.Name
	user.PhoneNumber = userChange.PhoneNumber
	user.Biography = userChange.Biography
	user.Birhtday = userChange.Birhtday
	user.Website = userChange.Website
	user.Gender = userChange.Gender
	handler.Service.UpdateUser(&user)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
}

func (handler *UserHandler) SendConfirmation(w http.ResponseWriter, r *http.Request) {
	b, _ := ioutil.ReadAll(r.Body)
	var user model.User
	user.Role = "user"
	err := json.Unmarshal(b, &user)
	if err != nil {
		var err model.Error
		err = model.SetError(err, "Error in reading payload.")
		json.NewEncoder(w).Encode(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var exists bool
	exists, err = handler.Service.UserExists(user.Email, user.Username)

	if exists == true {
		fmt.Printf("USAO u error ")
		var err model.Error
		err = model.SetError(err, "Already exists.")

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}

	handler.Service.SendEmail(user.Email)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
}

func (handler *UserHandler) SendEmailForAccountRecovery(w http.ResponseWriter, r *http.Request) {
	b, _ := ioutil.ReadAll(r.Body)
	var email EmailForRecovery
	err := json.Unmarshal(b, &email)
	fmt.Printf(email.Email)
	if err != nil {
		var err model.Error
		err = model.SetError(err, "Error in reading payload.")
		json.NewEncoder(w).Encode(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var exists bool
	exists, err = handler.Service.GetUserByEmail(email.Email)

	if exists == false {
		var err model.Error
		err = model.SetError(err, "User with that email doesn't exist.")

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}

	handler.Service.SendEmailForAccountRecovery(email.Email)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
}

func (handler *UserHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	b, _ := ioutil.ReadAll(r.Body)
	fmt.Printf(string(b))
	var user model.User
	user.Role = "user"
	err := json.Unmarshal(b, &user)
	if err != nil {
		var err model.Error
		err = model.SetError(err, "Error in reading payload.")
		json.NewEncoder(w).Encode(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = handler.Service.CreateUser(&user)
	if err != nil {
		var err model.Error
		err = model.SetError(err, "Failed in creating user.")
		json.NewEncoder(w).Encode(err)
		w.WriteHeader(http.StatusExpectationFailed)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (handler *UserHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	b, _ := ioutil.ReadAll(r.Body)
	var authDetails model.Authentication
	err := json.Unmarshal(b, &authDetails)
	if err != nil {
		var err model.Error
		err = model.SetError(err, "Error in reading payload.")

		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(err)
		return
	}
	var authUser model.User
	authUser = handler.Service.UserForLogin(authDetails.Email)

	if authUser.Email == "" {
		var err model.Error
		err = model.SetError(err, "Username or Password is incorrect")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}

	check := handler.Service.CheckPasswordHash(authDetails.Password, authUser.Password)

	if !check {
		var err model.Error
		err = model.SetError(err, "Username or Password is incorrect")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}

	validToken, err := handler.Service.GenerateJWT(authUser.Email, authUser.Role)
	if err != nil {
		var err model.Error
		err = model.SetError(err, "Failed to generate token")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}
	var token model.Token
	token.Email = authUser.Email
	token.Role = authUser.Role
	token.TokenString = validToken
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(token)
}

func (handler *UserHandler) GetUserByEmailAddress(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	email := vars["email"]

	user := handler.Service.GetUserByEmailAddress(email)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
