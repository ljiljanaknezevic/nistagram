package handler

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"search-service/model"
	"search-service/service"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	logrus "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

type SearchHandler struct {
	Service *service.SearchService
}

var log = logrus.New()

type FileWithBASE64 struct {
	Path     string `json:"path"`
	FileType string `json:"type"`
}

func init() {
	absPath, err := os.Getwd()

	path := filepath.Join(absPath, "files", "search-service.log")
	filel, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.Out = filel
	} else {
		log.WithFields(
			logrus.Fields{
				"location": "search-service.handler.searchHandler.init()",
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
			"location": "search-service.handler.searchHandler.init()",
		},
	).Info("Search-service Log file created/opened")
}
func (handler *SearchHandler) SaveComment(w http.ResponseWriter, r *http.Request) {

	var comment model.Comment
	comment.Text = r.PostFormValue("text")
	comment.PostID = r.PostFormValue("postID")
	comment.Email = r.PostFormValue("email")
	handler.Service.SaveComment(&comment)
	jsonResponse(w, http.StatusCreated, "File uploaded successfully!.")
}
func jsonResponse(w http.ResponseWriter, code int, message string) {
	//	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	fmt.Fprint(w, message)
}
func (handler *SearchHandler) GetUserByUsername(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]
	loggingUsername := vars["loggingUsername"]
	users := handler.Service.GetAllUsersExceptLogging(loggingUsername)
	loggedUser := handler.Service.GetUserByEmailAddress(loggingUsername)

	var result []model.User
	for _, element := range users {
		isBlocked := false
		amBlocked := false
		if element.Role == "user" {
			if strings.Contains(strings.ToLower(element.Username), strings.ToLower(username)) {
				//da li se u mojim blokovanim nalazi taj user
				if len(loggedUser.Blocked) != 0 {
					for _, elem := range loggedUser.Blocked {
						fmt.Println("//////////////////")
						fmt.Println(elem.Username)
						if elem.Username == element.Email {
							isBlocked = true
						}
					}
				}
				if len(loggedUser.UsersWhoBlocked) != 0 {
					for _, elem := range loggedUser.UsersWhoBlocked {
						fmt.Println(elem.Username)
						if elem.Username == element.Email {
							amBlocked = true
						}
					}
				}
				//ako ga nisam blokirala i nisam blokirana dodaj ga
				if !amBlocked && !isBlocked {
					result = append(result, element)
				}
			}
		}
	}

	json.NewEncoder(w).Encode(result)
}

func (handler *SearchHandler) GetUserByUsernameForUnregistredUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]

	users := handler.Service.GetAllUsers()
	var result []model.User

	for _, element := range users {
		if !element.IsPrivate && element.Role == "user" {
			if strings.Contains(strings.ToLower(element.Username), strings.ToLower(username)) {
				result = append(result, element)
			}

		}
	}

	json.NewEncoder(w).Encode(result)
}
func contains(s []model.Post, str model.Post) bool {
	for _, v := range s {
		if v.ID == str.ID {
			return true
		}
	}

	return false
}
func (handler *SearchHandler) SearchPostsByLocation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	location := vars["location"]
	email := vars["email"]

	loggedUser := handler.Service.GetUserByEmailAddress(email)
	posts := handler.Service.GetAllPosts()
	var result []model.Post

	for _, element := range posts {
		isBlocked := false
		amBlocked := false
		if element.Email != email {
			if !handler.Service.GetUserByEmailAddress(element.Email).IsPrivate {
				if strings.Contains(strings.ToLower(element.Location), strings.ToLower(location)) {
					if len(loggedUser.Blocked) != 0 {
						for _, elem := range loggedUser.Blocked {
							if elem.Username == element.Email {
								isBlocked = true
							}
						}
					}
					if len(loggedUser.UsersWhoBlocked) != 0 {
						for _, elem := range loggedUser.UsersWhoBlocked {
							fmt.Println(elem.Username)
							if elem.Username == element.Email {
								amBlocked = true
							}
						}
					}
					//ako ga nisam blokirala i nisam blokirana dodaj ga
					if !amBlocked && !isBlocked {
						result = append(result, element)
					}
				}
			} else {
				for _, follower := range handler.Service.GetUserByEmailAddress(element.Email).Followers {
					if strings.Compare(follower.Username, handler.Service.GetUserByEmailAddress(email).Username) == 0 {
						if strings.Contains(strings.ToLower(element.Location), strings.ToLower(location)) {
							if !contains(result, element) {
								result = append(result, element)
							}
						}
					}
				}
			}
		}
	}
	var newPosts []model.Post
	for _, element := range result {
		var s string = strconv.FormatUint(uint64(element.ID), 10)

		element.Comments = handler.Service.GetAllCommentsByPostsID(s)
		newPosts = append(newPosts, element)
	}

	json.NewEncoder(w).Encode(newPosts)
}

func (handler *SearchHandler) GetPostsForFeed(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	email := vars["email"]
	loggedUser := handler.Service.GetUserByEmailAddress(email)
	posts := handler.Service.GetAllPosts()
	var result []model.Post
	for _, elemet := range loggedUser.Following {
		for _, el := range posts {
			if handler.Service.GetUserByUsername(elemet.Username).Email == el.Email {
				result = append(result, el)
			}
		}
	}
	var res []model.Post
	var pom bool
	pom = false

	for _, elee := range result {
		isBlocked := false
		amBlocked := false
		pom = false
		for _, el := range loggedUser.Muted {
			if handler.Service.GetUserByUsername(el.Username).Email == elee.Email {
				pom = true
				break
			}
		}
		for _, elem := range loggedUser.Blocked {
			if elem.Username == elee.Email {
				isBlocked = true
			}
		}
		for _, elemm := range loggedUser.UsersWhoBlocked {
			if elemm.Username == elee.Email {
				amBlocked = true
			}
		}
		if !pom && !isBlocked && !amBlocked {
			res = append(res, elee)
		}
	}
	if len(res) >= 0 {
		var newPosts []model.Post
		for _, element := range res {
			var s string = strconv.FormatUint(uint64(element.ID), 10)
			element.Comments = handler.Service.GetAllCommentsByPostsID(s)
			newPosts = append(newPosts, element)
		}
		json.NewEncoder(w).Encode(newPosts)
	} else {
		var newPosts []model.Post
		for _, element := range result {
			var s string = strconv.FormatUint(uint64(element.ID), 10)
			element.Comments = handler.Service.GetAllCommentsByPostsID(s)
			newPosts = append(newPosts, element)
		}
		json.NewEncoder(w).Encode(newPosts)
	}

}

func (handler *SearchHandler) GetStoriesForFeed(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	email := vars["email"]
	loggedUser := handler.Service.GetUserByEmailAddress(email)
	stories := handler.Service.GetAllStories()
	var result []model.Story
	for _, elemet := range loggedUser.Following {
		for _, el := range stories {
			if handler.Service.GetUserByUsername(elemet.Username).Email == el.Email && el.CreatedAt.Add(time.Hour*time.Duration(24)).After(time.Now().Local()) {
				result = append(result, el)
			}
		}
	}

	var res []model.Story
	var pom bool
	pom = false

	for _, elee := range result {
		pom = false
		isBlocked := false
		amBlocked := false
		for _, el := range loggedUser.Muted {
			if handler.Service.GetUserByUsername(el.Username).Email == elee.Email {
				pom = true
				break
			}

		}
		for _, elem := range loggedUser.Blocked {
			if elem.Username == elee.Email {
				isBlocked = true
			}
		}
		for _, elemm := range loggedUser.UsersWhoBlocked {
			if elemm.Username == elee.Email {
				amBlocked = true
			}
		}
		if !pom && !amBlocked && !isBlocked {
			res = append(res, elee)
		}
	}
	if len(res) >= 0 {
		json.NewEncoder(w).Encode(res)
	} else {
		json.NewEncoder(w).Encode(result)
	}
}

func (handler *SearchHandler) SearchPostsByTag(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tag := vars["tag"]
	email := vars["email"]

	loggedUser := handler.Service.GetUserByEmailAddress(email)
	posts := handler.Service.GetAllPosts()
	var result []model.Post

	for _, element := range posts {
		isBlocked := false
		amBlocked := false
		if element.Email != email {
			if !handler.Service.GetUserByEmailAddress(element.Email).IsPrivate {
				if strings.Contains(strings.ToLower(element.Tags), strings.ToLower(tag)) {
					if len(loggedUser.Blocked) != 0 {
						for _, elem := range loggedUser.Blocked {
							if elem.Username == element.Email {
								isBlocked = true
							}
						}
					}
					if len(loggedUser.UsersWhoBlocked) != 0 {
						for _, elem := range loggedUser.UsersWhoBlocked {
							if elem.Username == element.Email {
								amBlocked = true
							}
						}
					}
					//ako ga nisam blokirala i nisam blokirana dodaj ga
					if !amBlocked && !isBlocked {
						result = append(result, element)
					}
				}
			} else {
				for _, follower := range handler.Service.GetUserByEmailAddress(element.Email).Followers {
					if strings.Compare(follower.Username, handler.Service.GetUserByEmailAddress(email).Username) == 0 {
						if strings.Contains(strings.ToLower(element.Tags), strings.ToLower(tag)) {
							if !contains(result, element) {
								result = append(result, element)
							}
						}
					}
				}
			}
		}
	}
	var newPosts []model.Post
	for _, element := range result {
		var s string = strconv.FormatUint(uint64(element.ID), 10)

		element.Comments = handler.Service.GetAllCommentsByPostsID(s)
		newPosts = append(newPosts, element)
	}
	json.NewEncoder(w).Encode(newPosts)
}

func (handler *SearchHandler) SearchPostsByLocationUnregistered(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	location := vars["location"]
	log.WithFields(logrus.Fields{
		"location": "search-service.handler.searchHandler.SearchPostsByLocationUnregistered()"}).Info("Search posts by location from unregistred user.")

	posts := handler.Service.GetAllPosts()
	var result []model.Post

	for _, element := range posts {
		if !handler.Service.GetUserByEmailAddress(element.Email).IsPrivate {
			if strings.Contains(strings.ToLower(element.Location), strings.ToLower(location)) {
				result = append(result, element)
			}
		}
	}
	var newPosts []model.Post
	for _, element := range result {
		var s string = strconv.FormatUint(uint64(element.ID), 10)

		element.Comments = handler.Service.GetAllCommentsByPostsID(s)
		newPosts = append(newPosts, element)

	}

	json.NewEncoder(w).Encode(newPosts)
}
func (handler *SearchHandler) SearchPostsByTagUnregistered(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tag := vars["tag"]

	posts := handler.Service.GetAllPosts()
	var result []model.Post

	for _, element := range posts {
		if !handler.Service.GetUserByEmailAddress(element.Email).IsPrivate {
			if strings.Contains(strings.ToLower(element.Tags), strings.ToLower(tag)) {
				result = append(result, element)
			}
		}

	}
	var newPosts []model.Post
	for _, element := range result {
		var s string = strconv.FormatUint(uint64(element.ID), 10)

		element.Comments = handler.Service.GetAllCommentsByPostsID(s)
		newPosts = append(newPosts, element)
	}

	json.NewEncoder(w).Encode(newPosts)
}

func (handler *SearchHandler) GetPostsForSearchedUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	email := vars["id"]
	emailLoggedUser := vars["email"]

	posts := handler.Service.GetPostsForSearchedUser(email)
	var newPosts []model.Post
	for _, element := range posts {
		var s string = strconv.FormatUint(uint64(element.ID), 10)

		element.Comments = handler.Service.GetAllCommentsByPostsID(s)
		newPosts = append(newPosts, element)
	}
	var user model.User
	user = handler.Service.GetUserByEmailAddress(emailLoggedUser)

	var searchUser model.User
	searchUser = handler.Service.GetUserByEmailAddress(email)

	if !searchUser.IsPrivate {

		fmt.Println(user.Username)
		json.NewEncoder(w).Encode(newPosts)
		w.WriteHeader(http.StatusOK)
		return
	} else {
		for _, following := range handler.Service.GetUserByEmailAddress(user.Email).Following {
			if strings.Compare(following.Username, handler.Service.GetUserByEmailAddress(email).Username) == 0 {
				json.NewEncoder(w).Encode(newPosts)
				w.WriteHeader(http.StatusOK)
				return
			}
		}
	}

	w.WriteHeader(http.StatusBadRequest)

}
func (handler *SearchHandler) GetPostsForSearchedUserUnregistered(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	email := vars["id"]
	posts := handler.Service.GetPostsForSearchedUser(email)
	var newPosts []model.Post
	for _, element := range posts {
		var s string = strconv.FormatUint(uint64(element.ID), 10)

		element.Comments = handler.Service.GetAllCommentsByPostsID(s)
		newPosts = append(newPosts, element)
	}
	json.NewEncoder(w).Encode(newPosts)

}
func (handler *SearchHandler) MediaForFront(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idString := vars["id"]
	id, _ := strconv.ParseUint(idString, 10, 64)
	log.WithFields(logrus.Fields{
		"location": "search-service.handler.searchHandler.MediaForFront()"}).Info("Get image by imageID.")

	var file model.File
	file = handler.Service.FindFileById(uint(id))

	buffer := new(bytes.Buffer)
	f, _ := os.Open(file.Path)

	defer f.Close()
	io.Copy(buffer, f)
	s := base64.StdEncoding.EncodeToString(buffer.Bytes())
	var base64String FileWithBASE64
	base64String.FileType = file.Type
	base64String.Path = s
	imagesMarshaled, err := json.Marshal(base64String)
	if err != nil {
		log.WithFields(logrus.Fields{
			"location": "search-service.handler.searchHandler.MediaForFront()"}).Error(err)

		fmt.Fprint(w, err)
	}
	log.WithFields(logrus.Fields{
		"location": "search-service.handler.searchHandler.MediaForFront()"}).Info("Get image by imageID success.")

	/*	image, _, _ := image.Decode(f)

		buffer := new(bytes.Buffer)
		if err := jpeg.Encode(buffer, image, nil); err != nil {
			if err := png.Encode(buffer, image); err != nil {
				log.WithFields(logrus.Fields{
					"location": "search-service.handler.searchHandler.MediaForFront()"}).Error("Unable to encode image.")

				fmt.Printf("Unable to encode image")
			}
		}

		mediaZaFront = buffer.Bytes()

		imagesMarshaled, err := json.Marshal(mediaZaFront)

		if err != nil {
			log.WithFields(logrus.Fields{
				"location": "search-service.handler.searchHandler.MediaForFront()"}).Error(err)

			fmt.Fprint(w, err)
		}
		log.WithFields(logrus.Fields{
			"location": "search-service.handler.searchHandler.MediaForFront()"}).Info("Get image by imageID success.")
	*/
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(imagesMarshaled)

	//json.NewEncoder(w).Encode(imagesMarshaled)

}

func (handler *SearchHandler) VideoZaFront(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idString := vars["id"]
	id, _ := strconv.ParseUint(idString, 10, 64)
	log.WithFields(logrus.Fields{
		"location": "search-service.handler.searchHandler.MediaForFront()"}).Info("Get image by imageID.")

	var file model.File
	file = handler.Service.FindFileById(uint(id))

	buffer := new(bytes.Buffer)
	f, _ := os.Open(file.Path)

	defer f.Close()
	io.Copy(buffer, f)
	//s:=string(buffer.Bytes())
	s := base64.StdEncoding.EncodeToString(buffer.Bytes())
	imagesMarshaled, err := json.Marshal(s)

	if err != nil {
		log.WithFields(logrus.Fields{
			"location": "search-service.handler.searchHandler.MediaForFront()"}).Error(err)

		fmt.Fprint(w, err)
	}
	log.WithFields(logrus.Fields{
		"location": "search-service.handler.searchHandler.MediaForFront()"}).Info("Get image by imageID success.")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(imagesMarshaled)

}
