package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/mail"
	"strings"
	"user-service-mod/model"
	"user-service-mod/service"

	"github.com/gorilla/mux"
	"regexp"
)

type UserHandler struct {
	Service *service.UserService
}
var  secretBase32 string
var authUser model.User

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
	if validEmail(changePassword.Email) && CheckPasswordLever(changePassword.Password)==nil && changePassword.Password==changePassword.ConfPassword{
	var user model.User
	user = handler.Service.GetUserByEmailAddress(changePassword.Email)
	if user.Username == "" {
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
	}else{
		var err model.Error
		err = model.SetError(err, "Incorrectly entered data.")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}
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
	if user.Username == "" {
		var err model.Error
		err = model.SetError(err, "User with that email doesn't exist.")

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}
	if validateUsername(userChange.Username){
		user.Email = userChange.Email
		user.Username = userChange.Username
		user.Name = userChange.Name
		user.PhoneNumber = userChange.PhoneNumber
		user.Biography = userChange.Biography
		user.Birhtday = userChange.Birhtday
		user.Website = userChange.Website
		user.Gender = userChange.Gender
		user.IsPrivate = userChange.IsPrivate
		handler.Service.UpdateUser(&user)
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		}else{
		var err model.Error
		err = model.SetError(err, "Incorrectly entered data.")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}
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
	if validateName(user.Name) && validateUsername(user.Username) && validEmail(user.Email) &&  CheckPasswordLever(user.Password) ==nil  {

		isValid := handler.Service.CreateUser(&user)
		if !isValid {
			var err model.Error
			err = model.SetError(err, "Failed in creating user.")
			w.WriteHeader(http.StatusExpectationFailed)
			json.NewEncoder(w).Encode(err)
			return
		}
		json.NewEncoder(w).Encode(user)
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")

	} else {
		var err model.Error
		err = model.SetError(err, "Incorrectly entered data.")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)

		return
	}
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
	if validEmail(authDetails.Email) && CheckPasswordLever(authDetails.Password)==nil{
		//var authUser model.User
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
		handler.Service.SendEmailWithQR(authUser.Email)
		w.WriteHeader(http.StatusOK)

	} else{
		var err model.Error
		err = model.SetError(err, "Incorrectly entered data.")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return

	}

}

func (handler *UserHandler) GetUserByEmailAddress(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	email := vars["email"]

	user := handler.Service.GetUserByEmailAddress(email)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (handler *UserHandler) Follow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	email := vars["email"]                       //maja
	followerUsername := vars["followerUsername"] //dragana
	var followerUser model.User
	followerUser = handler.Service.GetUserByUsername(followerUsername) //dragana
	var user model.User
	user = handler.Service.GetUserByEmailAddress(email)

	if followerUser.IsPrivate {
		var waitingFollower model.WaitingFollower
		waitingFollower.Username = user.Username
		followerUser.WaitingFollowers = append(followerUser.WaitingFollowers, waitingFollower)
		handler.Service.UpdateUser(&followerUser)
	} else {
		var follower model.Follower
		follower.Username = user.Username
		followerUser.Followers = append(followerUser.Followers, follower)
		handler.Service.UpdateUser(&followerUser)

		var following model.Following
		following.Username = followerUsername
		user.Following = append(user.Following, following)
		handler.Service.UpdateUser(&user)
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
func (handler *UserHandler) AcceptRequest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	email := vars["email"]
	followerUsername := vars["followerUsername"]
	var requestingUser model.User
	requestingUser = handler.Service.GetUserByUsername(followerUsername)
	var user model.User
	user = handler.Service.GetUserByEmailAddress(email)

	var following model.Following
	following.Username = user.Username
	requestingUser.Following = append(requestingUser.Following, following)
	handler.Service.UpdateUser(&requestingUser)

	var follower model.Follower
	follower.Username = followerUsername
	user.Followers = append(user.Followers, follower)
	handler.Service.UpdateUser(&user)

	var waitingFollower model.WaitingFollower
	waitingFollower.Username = followerUsername
	for _, element := range user.WaitingFollowers {
		if element.Username == followerUsername {
			handler.Service.DeleteFromWaitingList(element.ID)
		}
	}
	handler.Service.UpdateUser(&user)
	w.WriteHeader(http.StatusOK)
}

func (handler *UserHandler) DeclineRequest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	email := vars["email"]
	followerUsername := vars["followerUsername"]
	var user model.User
	user = handler.Service.GetUserByEmailAddress(email)

	var waitingFollower model.WaitingFollower
	waitingFollower.Username = followerUsername
	for _, element := range user.WaitingFollowers {
		if element.Username == followerUsername {
			handler.Service.DeleteFromWaitingList(element.ID)
		}
	}
	handler.Service.UpdateUser(&user)
	w.WriteHeader(http.StatusOK)
}

func RemoveIndex(s []model.WaitingFollower, index int) []model.WaitingFollower {
	return append(s[:index], s[index+1:]...)
}

func (handler *UserHandler) AlreadyFollow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	email := vars["email"]
	followerUsername := vars["followerUsername"]
	var user model.User
	user = handler.Service.GetUserByEmailAddress(email)
	var waitingUser model.User
	waitingUser = handler.Service.GetUserByUsername(followerUsername)
	for _, element := range user.Following {
		if strings.Compare(element.Username, followerUsername) == 0 {
			w.WriteHeader(http.StatusBadRequest)
		}

	}
	for _, element := range waitingUser.WaitingFollowers {
		if strings.Compare(element.Username, user.Username) == 0 {
			w.WriteHeader(http.StatusBadRequest)
		}

	}
	w.WriteHeader(http.StatusOK)
}
func (handler *UserHandler) GetAllFromWaitingList(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	email := vars["email"]
	var user model.User
	user = handler.Service.GetUserByEmailAddress(email)
	var result []model.User
	for _, elem := range user.WaitingFollowers {
		result = append(result, handler.Service.GetUserByUsername(elem.Username))
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(result)
}

func (handler *UserHandler) GetAllFollowers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	email := vars["email"]
	var user model.User
	user = handler.Service.GetUserByEmailAddress(email)
	var result []model.User
	for _, elem := range user.Followers {
		result = append(result, handler.Service.GetUserByUsername(elem.Username))
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(result)
}
func (handler *UserHandler) GetAllUsersExceptLogging(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	email := vars["email"]
	var users []model.User
	users = handler.Service.GetAllUsersExceptLogging(email)
	json.NewEncoder(w).Encode(users)
}

func validEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func validPassword(password string) bool {
	var strongRegex ="^(?=.*[a-z])(?=.*[A-Z])(?=.*[0-9])(?=.*[!@#\\$%\\^&\\*])(?=.{8,})"
	isValid , _ := regexp.MatchString(strongRegex, password)
	return isValid
}

func validateName(name string) bool{
	var pattern = "^[a-zA-Z]+[a-zA-Z\\s]*$"
	isValid, _ := regexp.MatchString(pattern, name)
	return isValid
}

func validateUsername(name string) bool{
	var pattern = "^[a-zA-Z0-9]+$"
	isValid, _ := regexp.MatchString(pattern, name)
	return isValid
}
func CheckPasswordLever(ps string) error {
	if len(ps) < 8 {
		return fmt.Errorf("password len is < 8")
	}
	num := `[0-9]{1}`
	a_z := `[a-z]{1}`
	A_Z := `[A-Z]{1}`
	symbol := `[!@#~$%^&*()+|_]{1}`
	if b, err := regexp.MatchString(num, ps); !b || err != nil {
		return fmt.Errorf("password need num :%v", err)
	}
	if b, err := regexp.MatchString(a_z, ps); !b || err != nil {
		return fmt.Errorf("password need a_z :%v", err)
	}
	if b, err := regexp.MatchString(A_Z, ps); !b || err != nil {
		return fmt.Errorf("password need A_Z :%v", err)
	}
	if b, err := regexp.MatchString(symbol, ps); !b || err != nil {
		return fmt.Errorf("password need symbol :%v", err)
	}
	return nil
}

func  (handler *UserHandler) HandlerFuncValidate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	input := vars["input"]

	if handler.Service.ValidateToken(input){
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
		handler.Service.SendEmailWithQR(authUser.Email)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(token)
		fmt.Println("Authorized")
		w.WriteHeader(http.StatusOK)
	}else{
		fmt.Println("Not authorized")
		w.WriteHeader(http.StatusBadRequest)
	}


}

