package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"user-service-mod/handler"
	"user-service-mod/repository"
	"user-service-mod/service"

	"user-service-mod/model"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

//--------GLOBAL VARIABLES---------------

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
		log.Fatalln("Invalid database url")
	}
	sqldb := connection.DB()

	err = sqldb.Ping()
	if err != nil {
		log.Fatal("Database connected")
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
}

//closes database connection
func CloseDatabase(connection *gorm.DB) {
	sqldb := connection.DB()
	sqldb.Close()
}

//check whether user is authorized or not
func IsAuthorized(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Header["Token"] == nil {
			var err model.Error
			err = model.SetError(err, "No Token Found")
			json.NewEncoder(w).Encode(err)
			return
		}

		var mySigningKey = []byte(secretkey)

		token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("There was an error in parsing token.")
			}
			return mySigningKey, nil
		})

		if err != nil {
			var err model.Error
			err = model.SetError(err, "Your Token has been expired.")
			json.NewEncoder(w).Encode(err)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			if claims["role"] == "admin" {
				r.Header.Set("Role", "admin")
				handler.ServeHTTP(w, r)
				return

			} else if claims["role"] == "user" {
				r.Header.Set("Role", "user")
				handler.ServeHTTP(w, r)
				return
			}
		}
		var reserr model.Error
		reserr = model.SetError(reserr, "Not Authorized.")
		json.NewEncoder(w).Encode(err)
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
	router.HandleFunc("/getAllRequests/{email}", handler.GetAllFromWaitingList).Methods("GET")
	router.HandleFunc("/getByEmail/{email}", handler.GetUserByEmailAddress).Methods("GET")
	router.HandleFunc("/changeUserData", handler.ChangeUserData).Methods("POST")
	router.HandleFunc("/follow/{followerUsername}/{email}", handler.Follow).Methods("POST")
	//router.HandleFunc("/declineRequest/{followerUsername}/{email}", handler.DeclineRequest).Methods("POST")
	router.HandleFunc("/acceptRequest/{followerUsername}/{email}", handler.AcceptRequest).Methods("POST")
	router.HandleFunc("/alreadyFollow/{followerUsername}/{email}", handler.AlreadyFollow).Methods("GET")
	/*router.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Access-Control-Request-Headers, Access-Control-Request-Method, Connection, Host, Origin, User-Agent, Referer, Cache-Control, X-header")
	})*/
}

//start the server
func ServerStart() {
	fmt.Println("Server started at http://localhost:8081")
	err := http.ListenAndServe(":8081", handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Access-Control-Allow-Origin", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(router))
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	db := GetDatabase()
	InitialMigration()
	CreateRouter()
	repo := initRepo(db)
	service := initServices(repo)
	handler := initHandler(service)
	InitializeRoute(handler)
	ServerStart()
}
