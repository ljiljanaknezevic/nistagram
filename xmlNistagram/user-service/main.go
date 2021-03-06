package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"strings"
	"user-service-mod/handler"
	"user-service-mod/repository"
	"user-service-mod/service"

	"user-service-mod/model"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	logrus "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

//--------GLOBAL VARIABLES---------------
var log = logrus.New()
var (
	router    *mux.Router
	secretkey string = "secretkeyjwt"
)

func initRepo(database *gorm.DB) *repository.UserRepository {
	return &repository.UserRepository{Database: database}
}

func initServices(repo *repository.UserRepository) *service.UserService {
	return &service.UserService{Repo: repo}
}

func initHandler(service *service.UserService) *handler.UserHandler {
	return &handler.UserHandler{Service: service}
}

//-------------DATABASE FUNCTIONS---------------------

//returns database connection
func GetDatabase() *gorm.DB {
	databasename := "postgres"
	database := "postgres"
	databasepassword := "super"
	databaseurl := "postgres://postgres:" + databasepassword + "@postgresdb:5432/" + databasename + "?sslmode=disable"
	//databaseurl := "host=localhost port=5432 user=postgres password=super dbname=postgres sslmode=disable"
	connection, err := gorm.Open(database, databaseurl)
	if err != nil {
		log.WithFields(logrus.Fields{
			"location": "user-service.main.GetDatabase()"}).Fatal("Invalid database url")

		log.Fatalln("Invalid database url")
	}
	sqldb := connection.DB()

	err = sqldb.Ping()
	if err != nil {
		log.WithFields(logrus.Fields{
			"location": "user-service.main.GetDatabase()"}).Fatal("User-service Database connected")

		//	log.Fatal("Database connected")
	}
	fmt.Println("Database connection successful.")
	return connection
}

//create user table in userdb
func InitialMigration() {
	connection := GetDatabase()
	defer CloseDatabase(connection)
	connection.AutoMigrate(model.User{})
	connection.AutoMigrate(model.Follower{})
	connection.AutoMigrate(model.WaitingFollower{})
	connection.AutoMigrate(model.Following{})
	connection.AutoMigrate(model.VerificationRequest{})
	connection.AutoMigrate(model.Blocked{})
	connection.AutoMigrate(model.Muted{})
	connection.AutoMigrate(model.UsersWhoBlocked{})
}

//closes database connection
func CloseDatabase(connection *gorm.DB) {
	sqldb := connection.DB()
	sqldb.Close()
}

//check whether user is authorized or not
func IsAuthorized(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.WithFields(
			logrus.Fields{"location": "user-service.main.IsAuthorized()"}).Info("Check is user authorized")

		if r.Header["Authorization"] == nil {
			log.WithFields(
				logrus.Fields{
					"location": "user-service.main.IsAuthorized()",
				},
			).Error("No Token Found")
			var err model.Error
			err = model.SetError(err, "No Token Found")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(err)
			return
		}

		var mySigningKey = []byte(secretkey)
		token, err := jwt.Parse(strings.Split(r.Header["Authorization"][0], " ")[1], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				log.WithFields(
					logrus.Fields{"location": "user-service.main.IsAuthorized()"}).Error("There was an error in parsing token.")
				return nil, fmt.Errorf("There was an error in parsing token.")
			}
			return mySigningKey, nil
		})

		if err != nil {
			log.WithFields(
				logrus.Fields{"location": "user-service.main.IsAuthorized()"}).Error("Your Token has been expired.")

			var err model.Error
			err = model.SetError(err, "Your Token has been expired.")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(err)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

			if claims["role"] == "admin" {
				log.WithFields(
					logrus.Fields{"location": "user-service.main.IsAuthorized()"}).Info("Admin authorize success.")
				r.Header.Set("Role", "admin")
				handler.ServeHTTP(w, r)
				return

			} else if claims["role"] == "user" {
				log.WithFields(
					logrus.Fields{"location": "user-service.main.IsAuthorized()"}).Info("User authorize success.")
				r.Header.Set("Role", "user")

				handler.ServeHTTP(w, r)
				return
			}
		}
		log.WithFields(
			logrus.Fields{"location": "user-service.main.IsAuthorized()"}).Error("User authorize fail.")

		var reserr model.Error
		reserr = model.SetError(reserr, "Not Authorized.")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(reserr)
	}
}

//----------------------ROUTES-------------------------------
//create a mux router
func CreateRouter() {
	router = mux.NewRouter()
}

//initialize all routes
func InitializeRoute(handler *handler.UserHandler) {
	router.HandleFunc("/signup", handler.SignUp).Methods("POST")
	router.HandleFunc("/signin", handler.SignIn).Methods("POST")
	router.HandleFunc("/confirmRegistration", handler.SendConfirmation).Methods("POST")
	router.HandleFunc("/sendEmailForAccountRecovery", handler.SendEmailForAccountRecovery).Methods("POST")
	router.HandleFunc("/changePassword", handler.ChangePassword).Methods("POST")
	router.HandleFunc("/getAllRequests/{email}", IsAuthorized(handler.GetAllFromWaitingList)).Methods("GET")
	router.HandleFunc("/getByEmail/{email}", handler.GetUserByEmailAddress).Methods("GET")
	router.HandleFunc("/changeUserData", IsAuthorized(handler.ChangeUserData)).Methods("POST")
	router.HandleFunc("/follow/{followerUsername}/{email}", IsAuthorized(handler.Follow)).Methods("POST")
	router.HandleFunc("/block/{followerUsername}/{email}", IsAuthorized(handler.Block)).Methods("POST")
	router.HandleFunc("/declineRequest/{followerUsername}/{email}", IsAuthorized(handler.DeclineRequest)).Methods("POST")
	router.HandleFunc("/acceptRequest/{followerUsername}/{email}", IsAuthorized(handler.AcceptRequest)).Methods("POST")
	router.HandleFunc("/alreadyFollow/{followerUsername}/{email}", IsAuthorized(handler.AlreadyFollow)).Methods("GET")
	router.HandleFunc("/getAllFollowers/{email}", IsAuthorized(handler.GetAllFollowers)).Methods("GET")
	router.HandleFunc("/getAllUsersExceptLogging/{email}", IsAuthorized(handler.GetAllUsersExceptLogging)).Methods("GET")
	router.HandleFunc("/getAllUsersExceptLoggingForTag/{email}", IsAuthorized(handler.GetAllUsersExceptLoggingForTag)).Methods("GET")
	router.HandleFunc("/validateToken/{input}", handler.HandlerFuncValidate).Methods("GET")
	router.HandleFunc("/createRequest", handler.CreateRequest).Methods("POST")
	router.HandleFunc("/getAllRequests", handler.GetAllRequestes).Methods("GET")
	router.HandleFunc("/acceptVerification/{email}", handler.AcceptVerification).Methods("POST")
	router.HandleFunc("/declineVerification/{email}", handler.DeclineVerification).Methods("POST")

	router.HandleFunc("/muteAccount/{username}/{email}", handler.MuteAccount).Methods("POST")
	router.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Access-Control-Request-Headers,Token, Access-Control-Request-Method, Connection, Host, Origin, User-Agent, Referer, Cache-Control, X-header")
	})
}

//start the server
func ServerStart() {
	fmt.Println("Server started at http://localhost:8081")
	err := http.ListenAndServe(":8081", handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Access-Control-Allow-Origin", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(router))
	if err != nil {
		log.WithFields(logrus.Fields{
			"location": "user-service.main.ServerStart()"}).Fatal(err)

		//log.Fatal(err)
	}
}

func main() {

	absPath, err := os.Getwd()

	path := filepath.Join(absPath, "files", "user-service.log")
	filel, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.Out = filel
	} else {
		log.WithFields(
			logrus.Fields{
				"location": "user-service.handler.userHandler.init()",
			},
		).Info("Failed to log to file, using default stderr")
	}
	log.SetOutput(&lumberjack.Logger{
		Filename:   path,
		MaxSize:    300, // megabytes
		MaxBackups: 3,
		MaxAge:     28,   //days
		Compress:   true, // disabled by default
	})
	log.SetFormatter(&logrus.JSONFormatter{})
	log.WithFields(
		logrus.Fields{
			"location": "user-service.handler.userHandler.init()",
		},
	).Info("User-service Log file created/opened")

	db := GetDatabase()
	InitialMigration()
	CreateRouter()
	repo := initRepo(db)
	service := initServices(repo)
	handler := initHandler(service)
	InitializeRoute(handler)
	ServerStart()
}
