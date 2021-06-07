package main

import (
	"fmt"
	"log"
	"net/http"
	"search-service/repository"
	"search-service/handler"
	"search-service/model"
	"search-service/service"


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

func initRepo(database *gorm.DB) *repository.SearchRepository {
	return &repository.SearchRepository{Database: database}
}

func initServices(repo *repository.SearchRepository) *service.SearchService {
	return &service.SearchService{Repo: repo}
}

func initHandler(service *service.SearchService) *handler.SearchHandler {
	return &handler.SearchHandler{Service: service}
}

//-------------DATABASE FUNCTIONS---------------------

//returns database connection
func GetUserDatabase() *gorm.DB {
	databasename := "postgres"
	database := "postgres"
	databasepassword := "super"
	//databaseurl := "postgres://postgres:" + databasepassword + "@localhost:5433/" + databasename + "?sslmode=disable"
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
func InitialUserMigration() {
	connection := GetUserDatabase()
	defer CloseUserDatabase(connection)
	connection.AutoMigrate(model.User{})
	connection.AutoMigrate(model.Follower{})
	connection.AutoMigrate(model.WaitingFollower{})
	connection.AutoMigrate(model.Following{})
}

//closes database connection
func CloseUserDatabase(connection *gorm.DB) {
	sqldb := connection.DB()
	sqldb.Close()
}

//----------------------ROUTES-------------------------------
//create a mux router
func CreateRouter() {
	router = mux.NewRouter()
}

//initialize all routes
func InitializeRoute(handler *handler.SearchHandler) {
		router.HandleFunc("/searchUserByUsername/{username}/{loggingUsername}", handler.GetUserByUsername).Methods("GET")
		//router.HandleFunc("/getAllUsers", handler.GetAllUsers).Methods("GET")
		router.HandleFunc("/searchUserByUsernameForUnregistredUser/{username}", handler.GetUserByUsernameForUnregistredUser).Methods("GET")
		router.HandleFunc("/searchPostByLocation/{location}/{email}", handler.SearchPostsByLocation).Methods("GET")
		router.HandleFunc("/searchPostByLocationUnregistered/{location}", handler.SearchPostsByLocationUnregistered).Methods("GET")
		router.HandleFunc("/getPostsForSearchedUser/{id}/{email}", handler.GetPostsForSearchedUser).Methods("GET")
		router.HandleFunc("/getMedia/{id}", handler.MediaForFront).Methods("GET")

	router.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Access-Control-Request-Headers, Access-Control-Request-Method, Connection, Host, Origin, User-Agent, Referer, Cache-Control, X-header")
		})
}

//start the server
func ServerStart() {
	fmt.Println("Server started at http://localhost:8083")
	err := http.ListenAndServe(":8083", handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Access-Control-Allow-Origin", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(router))
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	db := GetUserDatabase()
	InitialUserMigration()
	CreateRouter()
	repo := initRepo(db)
	service := initServices(repo)
	handler := initHandler(service)
	InitializeRoute(handler)
	ServerStart()
}
