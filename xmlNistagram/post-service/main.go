package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"post-service-mod/handler"
	"post-service-mod/model"
	"post-service-mod/repository"
	"post-service-mod/service"
	"strings"

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

func initRepo(database *gorm.DB) *repository.PostRepository {
	return &repository.PostRepository{Database: database}
}

//comment repo
func initCommentRepo(database *gorm.DB) *repository.CommentRepository {
	return &repository.CommentRepository{Database: database}
}
func initFileRepo(database *gorm.DB) *repository.FileRepository {
	return &repository.FileRepository{Database: database}
}

func initServices(repo *repository.PostRepository, fileRepo *repository.FileRepository, commentRepo *repository.CommentRepository) *service.PostService {
	return &service.PostService{Repo: repo, FileRepo: fileRepo, CommentRepo: commentRepo}
}

func initHandler(service *service.PostService) *handler.PostHandler {
	return &handler.PostHandler{Service: service}
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
			"location": "post-service.main.GetDatabase()"}).Fatal("Invalid database url")

		log.Fatalln("Invalid database url")
	}
	sqldb := connection.DB()

	err = sqldb.Ping()
	if err != nil {
		log.WithFields(logrus.Fields{
			"location": "post-service.main.GetDatabase()"}).Fatal("Post-service Database connected")

		//log.Fatal("Database connected")
	}
	fmt.Println("Database connection successful.")
	return connection
}

//create user table in userdb
func InitialMigration() {
	connection := GetDatabase()
	defer CloseDatabase(connection)
	connection.AutoMigrate(model.Post{})
	connection.AutoMigrate(model.File{})
	connection.AutoMigrate(model.Like{})
	connection.AutoMigrate(model.Spam{})
	connection.AutoMigrate(model.Comment{})
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
			logrus.Fields{"location": "post-service.main.IsAuthorized()"}).Info("Check is user authorized")

		if r.Header["Authorization"] == nil {
			log.WithFields(
				logrus.Fields{
					"location": "post-service.main.IsAuthorized()",
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
					logrus.Fields{"location": "post-service.main.IsAuthorized()"}).Error("There was an error in parsing token.")
				return nil, fmt.Errorf("There was an error in parsing token.")
			}
			return mySigningKey, nil
		})

		if err != nil {
			log.WithFields(
				logrus.Fields{"location": "post-service.main.IsAuthorized()"}).Error("Your Token has been expired.")
			var err model.Error
			err = model.SetError(err, "Your Token has been expired.")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(err)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			if claims["role"] == "admin" {
				log.WithFields(
					logrus.Fields{"location": "post-service.main.IsAuthorized()"}).Info("User authorize success.")
				r.Header.Set("Role", "admin")
				handler.ServeHTTP(w, r)
				return

			} else if claims["role"] == "user" {
				log.WithFields(
					logrus.Fields{"location": "post-service.main.IsAuthorized()"}).Info("User authorize success.")
				r.Header.Set("Role", "user")

				handler.ServeHTTP(w, r)
				return
			}
		}

		log.WithFields(
			logrus.Fields{"location": "post-service.main.IsAuthorized()"}).Error("User authorize fail.")
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
func InitializeRoute(handler *handler.PostHandler) {
	router.HandleFunc("/savePost", IsAuthorized(handler.SavePost)).Methods("POST")
	router.HandleFunc("/getAllPostsByEmail/{email}", IsAuthorized(handler.GetAllPostsByEmail)).Methods("GET")
	router.HandleFunc("/getImageByImageID/{imageID}", IsAuthorized(handler.GetImageByImageID)).Methods("GET")
	router.HandleFunc("/liked/{postID}/{userWhoLiked}", IsAuthorized(handler.Liked)).Methods("POST")
	router.HandleFunc("/getAllLikedPostsByEmail/{email}", IsAuthorized(handler.GetAllLikedPostsByEmail)).Methods("GET")
	router.HandleFunc("/reportPost", IsAuthorized(handler.CreateSpam)).Methods("POST")
	router.HandleFunc("/saveComment", handler.SaveComment).Methods("POST")
	router.HandleFunc("/getAllCommentsByPostsID/{postID}", handler.GetAllCommentsByPostsID).Methods("GET")

	router.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Access-Control-Request-Headers, Access-Control-Request-Method, Connection, Host, Origin, User-Agent, Referer, Cache-Control, X-header")
	})
}

//start the server
func ServerStart() {
	fmt.Println("Server started at http://localhost:8084")
	err := http.ListenAndServe(":8084", handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Access-Control-Allow-Origin", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(router))
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	absPath, err := os.Getwd()

	path := filepath.Join(absPath, "files", "post-service.log")
	filel, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.Out = filel
	} else {
		log.Info("Failed to log to file, using default stderr")
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
			"location": "post-service.main.go",
		},
	).Info("Log file created")

	db := GetDatabase()
	InitialMigration()
	CreateRouter()
	repo := initRepo(db)
	commentRepo := initCommentRepo(db)
	fileRepo := initFileRepo(db)
	service := initServices(repo, fileRepo, commentRepo)
	handler := initHandler(service)
	InitializeRoute(handler)
	ServerStart()
}
