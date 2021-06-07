package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"image"
	"image/jpeg"
	"image/png"
	"net/http"
	"os"
	"search-service/model"
	"search-service/service"
	"strconv"
	"strings"
)

type SearchHandler struct {
	Service *service.SearchService
}

func (handler *SearchHandler) GetUserByUsername(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]
	loggingUsername := vars["loggingUsername"]

	users := handler.Service.GetAllUsersExceptLogging(loggingUsername)
	var result []model.User

	for _, element := range users {
	if(strings.Contains(strings.ToLower(element.Username),strings.ToLower(username))){
		result = append(result, element)
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
		if !element.IsPrivate {
			if (strings.Contains(strings.ToLower(element.Username), strings.ToLower(username))) {
				result = append(result, element)
			}
		}
	}
	json.NewEncoder(w).Encode(result)
}
func contains(s []model.Post, str model.Post) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
func (handler *SearchHandler) SearchPostsByLocation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	location := vars["location"]
	email := vars["email"]

	posts := handler.Service.GetAllPosts()
	var result []model.Post

	for _, element := range posts {
		if element.Email!=email{
			if !handler.Service.GetUserByEmailAddress(element.Email).IsPrivate {
				if strings.Contains(strings.ToLower(element.Location),strings.ToLower(location)) {
					result = append(result, element)
				}
			}else {
				for _,follower := range handler.Service.GetUserByEmailAddress(element.Email).Followers{
					if strings.Compare(follower.Username,handler.Service.GetUserByEmailAddress(email).Username)==0{
						if strings.Contains(strings.ToLower(element.Location),strings.ToLower(location)) {
							if !contains(result,element) {
								result = append(result, element)
							}
						}
					}
				}
			}
		}
	}
	json.NewEncoder(w).Encode(result)
}
func (handler *SearchHandler) SearchPostsByLocationUnregistered(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	location := vars["location"]
	email := vars["email"]

	posts := handler.Service.GetAllPosts()
	var result []model.Post

	for _, element := range posts {
		if element.Email!=email{
			if !handler.Service.GetUserByEmailAddress(element.Email).IsPrivate {
				if strings.Contains(strings.ToLower(element.Location),strings.ToLower(location)) {
					result = append(result, element)
				}
			}

		}
	}
	json.NewEncoder(w).Encode(result)
}

func (handler *SearchHandler) GetPostsForSearchedUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	email := vars["id"]
	emailLoggedUser := vars["email"]

	posts := handler.Service.GetPostsForSearchedUser(email)

	var user model.User
	user = handler.Service.GetUserByEmailAddress(emailLoggedUser)

	var searchUser model.User
	searchUser = handler.Service.GetUserByEmailAddress(email)

	if !searchUser.IsPrivate {
		fmt.Println(user.Username)
		json.NewEncoder(w).Encode(posts)
		w.WriteHeader(http.StatusOK)
		return
	} else{
		for _, following := range handler.Service.GetUserByEmailAddress(user.Email).Following {
			if strings.Compare(following.Username, handler.Service.GetUserByEmailAddress(email).Username) == 0 {
				json.NewEncoder(w).Encode(posts)
				w.WriteHeader(http.StatusOK)
				return
			}
		}
	}

	w.WriteHeader(http.StatusBadRequest)

}
func (handler *SearchHandler) MediaForFront(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idString := vars["id"]
	id, _ := strconv.ParseUint(idString, 10, 64)
	var file model.File
	file = handler.Service.FindFileById(uint(id))

	var mediaZaFront []byte

	f, _ := os.Open(file.Path)

	defer f.Close()

	image, _, _ := image.Decode(f)


	buffer := new(bytes.Buffer)
	if err := jpeg.Encode(buffer, image, nil) ; err != nil {
		if err :=png.Encode(buffer,image); err!=nil {
			fmt.Printf("Unable to encode image")
		}
	}

	mediaZaFront = buffer.Bytes()

	imagesMarshaled, err := json.Marshal(mediaZaFront)

	if err != nil {
		fmt.Fprint(w, err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(imagesMarshaled)

	//json.NewEncoder(w).Encode(imagesMarshaled)

}