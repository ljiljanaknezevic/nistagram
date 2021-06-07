package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"image"
	"image/jpeg"
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
		if(strings.Contains(strings.ToLower(element.Username),strings.ToLower(username))){
			result = append(result, element)
		}
	}
	json.NewEncoder(w).Encode(result)
}

func (handler *SearchHandler) SearchPostsByLocation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	location := vars["location"]

	posts := handler.Service.GetAllPosts()
	var result []model.Post

	for _, element := range posts {
		if strings.Contains(strings.ToLower(element.Location),strings.ToLower(location)) {
			result = append(result, element)
		}
	}
	json.NewEncoder(w).Encode(result)
}

func (handler *SearchHandler) MediaForFront(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idString := vars["id"]
	id, _ := strconv.ParseUint(idString, 10, 64)
	var file model.File
	file = handler.Service.FindFileById(uint(id))

	var mediaZaFront []byte

	f, _ := os.Open(file.Path)
	fmt.Println(f)
	defer f.Close()

	image, _, _ := image.Decode(f)

	buffer := new(bytes.Buffer)
	if err := jpeg.Encode(buffer, image, nil); err != nil {
		fmt.Printf("Unable to encode image")
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