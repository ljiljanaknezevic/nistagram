package service

import (
	"crypto/tls"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"time"
	"user-service-mod/model"
	"user-service-mod/repository"
	gomail "gopkg.in/mail.v2"
)

var (
	secretkey string = "secretkeyjwt"
)

type UserService struct {
	Repo *repository.UserRepository
}

func (service *UserService) CreateUser(user *model.User) error {
	user.Password, _ = service.GeneratehashPassword(user.Password)
	service.Repo.CreateUser(user)
	return nil
}

func (service *UserService) UpdateUser(user *model.User) error {
	service.Repo.UpdateUser(user)
	return nil
}

func (service *UserService) UserExists(email string, username string) (bool, error) {
	exists := service.Repo.UserExists(email, username)
	return exists, nil
}


func (service *UserService) GetUserByEmail(email string) (bool, error) {
	exists := service.Repo.GetUserByEmail(email)
	return exists, nil
}

func (service *UserService) GeneratehashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (service *UserService) UserForLogin(email string) model.User {
	user := service.Repo.UserForLogin(email)
	return user
}

func (service *UserService) GetUserByEmailAddress(email string) model.User {
	user := service.Repo.GetUserByEmailAddress(email)
	return user
}


func (service *UserService) CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

//Generate JWT token
func (service *UserService) GenerateJWT(email, role string) (string, error) {
	var mySigningKey = []byte(secretkey)
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims["email"] = email
	claims["role"] = role
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		fmt.Errorf("Something went Wrong: %s", err.Error())
		return "", err
	}

	return tokenString, nil
}

func (service *UserService) SendEmail(email string) {
	m := gomail.NewMessage()

	// Set E-Mail sender
	m.SetHeader("From", "notificationsnotifications22@gmail.com")

	// Set E-Mail receivers
	m.SetHeader("To", email)

	// Set E-Mail subject
	m.SetHeader("Subject", "Confirm registration")
	// Set E-Mail body. You can set plain text or html with text/html
	m.SetBody("text/html", "<a href='"+ "http://localhost:8082/confirmRegistration.html" + "'>Confirm registration!</a>")

	// Settings for SMTP server
	d := gomail.NewDialer("smtp.gmail.com", 587, "notificationsnotifications22@gmail.com", "Admin123#")

	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Now send E-Mail
	if err := d.DialAndSend(m); err != nil {
		fmt.Println(err)
		panic(err)
	}

	return
}

func (service *UserService) SendEmailForAccountRecovery(email string) {
	m := gomail.NewMessage()

	// Set E-Mail sender
	m.SetHeader("From", "notificationsnotifications22@gmail.com")

	// Set E-Mail receivers
	m.SetHeader("To", email)

	// Set E-Mail subject
	m.SetHeader("Subject", "Account recovery")
	// Set E-Mail body. You can set plain text or html with text/html
	m.SetBody("text/html", "<a href='"+ "http://localhost:8082/recoveryAccount.html" + "'>Change password</a>")

	// Settings for SMTP server
	d := gomail.NewDialer("smtp.gmail.com", 587, "notificationsnotifications22@gmail.com", "Admin123#")

	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Now send E-Mail
	if err := d.DialAndSend(m); err != nil {
		fmt.Println(err)
		panic(err)
	}

	return
}


