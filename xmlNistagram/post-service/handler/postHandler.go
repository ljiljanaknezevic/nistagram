package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"post-service-mod/model"
	"post-service-mod/service"
)

type PostHandler struct {
	Service *service.PostService
}

func (handler *PostHandler) SavePost(w http.ResponseWriter, r *http.Request) {

	r.ParseMultipartForm(32 << 20)
	file, handle, err := r.FormFile("file")

	if err != nil {
		fmt.Println(err)
		fmt.Fprintf(w, "%v", err)
		return
	}
	defer file.Close()
	absPath, err := os.Getwd()

	path := filepath.Join(absPath, "files", handle.Filename)
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0666)

	if err != nil {
		http.Error(w, "Expected file", http.StatusBadRequest)
		return
	}
	io.Copy(f, file)

	var savingFile model.File
	savingFile.Path = path
	savingFile.Type = "IMAGE"

	err = handler.Service.SaveFile(&savingFile)

	if err != nil {
		var err model.Error
		err = model.SetError(err, "Failed in creating file.")
		json.NewEncoder(w).Encode(err)
		w.WriteHeader(http.StatusExpectationFailed)
		return
	}

	fileId := handler.Service.FindFileIdByPath(path)

	fmt.Println(fileId)
	fmt.Println(r.PostFormValue("description"))
	fmt.Println(r.PostFormValue("location"))
	fmt.Println(r.PostFormValue("tags"))

	var post model.Post
	post.Description = r.PostFormValue("description")
	post.Location = r.PostFormValue("location")
	post.Tags = r.PostFormValue("tags")
	post.ImageID = fileId
	post.Email = r.PostFormValue("email")
	handler.Service.SavePost(&post)
	jsonResponse(w, http.StatusCreated, "File uploaded successfully!.")

}

func jsonResponse(w http.ResponseWriter, code int, message string) {
	//	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	fmt.Fprint(w, message)
}
