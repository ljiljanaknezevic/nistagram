package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"post-service-mod/model"
	"post-service-mod/service"
	"strconv"

	"github.com/gorilla/mux"
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

func (handler *PostHandler) GetAllPostsByEmail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	email := vars["email"]
	var result []model.Post
	result = handler.Service.GetAllPostsByEmail(email)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
func (handler *PostHandler) GetImageByImageID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	imageID := vars["imageID"]

	u64, err := strconv.ParseUint(imageID, 10, 32)
	if err != nil {
		fmt.Println(err)
	}

	var image_ID uint
	image_ID = uint(u64)
	var imagePath string
	imagePath = handler.Service.FindFilePathById(image_ID)

	//tamara
	var mediaZaFront []byte
	//image2 je putanja
	f, _ := os.Open(imagePath)
	defer f.Close()
	image, _, _ := image.Decode(f)
	buffer := new(bytes.Buffer)
	if err := jpeg.Encode(buffer, image, nil); err != nil {
		log.Println("unable to encode image.")
	}
	mediaZaFront = buffer.Bytes()
	imagesMarshaled, err := json.Marshal(mediaZaFront)
	if err != nil {
		fmt.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(imagesMarshaled)
}
