package handler

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"search-service/model"
	"search-service/service"
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
	if(strings.Contains(element.Username,username)){
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
		if(strings.Contains(element.Username,username)){
			result = append(result, element)
		}
	}
	json.NewEncoder(w).Encode(result)
}