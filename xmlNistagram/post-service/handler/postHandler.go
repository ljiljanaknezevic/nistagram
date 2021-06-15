package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"post-service-mod/model"
	"post-service-mod/service"
	"strconv"
	"text/template"

	"github.com/gorilla/mux"
	logrus "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

type PostHandler struct {
	Service *service.PostService
}

var log = logrus.New()

func init() {
	absPath, err := os.Getwd()

	path := filepath.Join(absPath, "files", "post-service.log")
	filel, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.Out = filel
	} else {
		log.WithFields(
			logrus.Fields{
				"location": "post-service.handler.postHandler.init()",
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
			"location": "post-service.handler.postHandler.init()",
		},
	).Info("Post-service Log file created/opened")
}
func (handler *PostHandler) SavePost(w http.ResponseWriter, r *http.Request) {
	log.WithFields(logrus.Fields{
		"location":   "post-service.handler.postHandler.SavePost()",
		"user_email": template.HTMLEscapeString(r.PostFormValue("email"))}).Info("User add post.")
	r.ParseMultipartForm(32 << 20)
	file, handle, err := r.FormFile("file")

	if err != nil {
		log.WithFields(
			logrus.Fields{
				"location": "post-service.handler.postHandler.SavePost()",
			},
		).Error(err)
		fmt.Println(err)
		fmt.Fprintf(w, "%v", err)
		return
	}
	defer file.Close()
	absPath, err := os.Getwd()
	var fileType string
	fileType = r.FormValue("type")

	path := filepath.Join(absPath, "files", handle.Filename)

	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0666)

	if err != nil {
		log.WithFields(
			logrus.Fields{
				"location": "post-service.handler.postHandler.SavePost()",
			},
		).Error("Bad format file")
		http.Error(w, "Expected file", http.StatusBadRequest)
		return
	}
	io.Copy(f, file)

	var savingFile model.File
	savingFile.Path = path
	savingFile.Type = fileType

	err = handler.Service.SaveFile(&savingFile)

	if err != nil {
		log.WithFields(logrus.Fields{
			"location": "post-service.handler.postHandler.SavePost()"}).Error("Failed in creating file.")
		var err model.Error
		err = model.SetError(err, "Failed in creating file.")
		json.NewEncoder(w).Encode(err)
		w.WriteHeader(http.StatusExpectationFailed)
		return
	}
	log.WithFields(logrus.Fields{
		"location": "post-service.handler.postHandler.SavePost()"}).Info("File localy created.")
	fileId := handler.Service.FindFileIdByPath(path)

	var post model.Post
	post.Description = r.PostFormValue("description")
	post.Location = r.PostFormValue("location")
	post.Tags = r.PostFormValue("tags")
	post.ImageID = fileId
	post.Email = r.PostFormValue("email")
	err = handler.Service.SavePost(&post)

	if err != nil {
		log.WithFields(logrus.Fields{
			"location": "post-service.handler.postHandler.SavePost()"}).Error("Save post failed.")
		var err model.Error
		err = model.SetError(err, "Save post fail.")
		json.NewEncoder(w).Encode(err)
		w.WriteHeader(http.StatusExpectationFailed)
		return
	}
	log.WithFields(logrus.Fields{
		"location": "post-service.handler.postHandler.SavePost()"}).Info("Save post success.")
	log.WithFields(logrus.Fields{
		"location":   "post-service.handler.postHandler.SavePost()",
		"user_email": template.HTMLEscapeString(r.PostFormValue("email"))}).Info("User add post success.")

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
	log.WithFields(logrus.Fields{
		"location":   "post-service.handler.postHandler.GetAllPostsByEmail()",
		"user_email": template.HTMLEscapeString(email)}).Info("Get all posts for user.")
	var result []model.Post
	result = handler.Service.GetAllPostsByEmail(email)

	if result == nil {
		log.WithFields(logrus.Fields{
			"location":   "post-service.handler.postHandler.GetAllPostsByEmail()",
			"user_email": template.HTMLEscapeString(email)}).Error("Get all posts for user fail.")
		var err model.Error
		err = model.SetError(err, "Get all posts fail.")
		json.NewEncoder(w).Encode(err)
		w.WriteHeader(http.StatusExpectationFailed)
		return
	}
	log.WithFields(logrus.Fields{
		"location":   "post-service.handler.postHandler.GetAllPostsByEmail()",
		"user_email": template.HTMLEscapeString(email)}).Info("Get all posts for user success.")
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
func (handler *PostHandler) GetImageByImageID(w http.ResponseWriter, r *http.Request) {
	log.WithFields(logrus.Fields{
		"location": "post-service.handler.postHandler.GetImageByImageID()"}).Info("Get image by imageID.")
	vars := mux.Vars(r)
	imageID := vars["imageID"]
	u64, err := strconv.ParseUint(imageID, 10, 32)
	if err != nil {
		log.WithFields(logrus.Fields{
			"location": "post-service.handler.postHandler.GetImageByImageID()"}).Error(err)
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

		log.WithFields(logrus.Fields{
			"location": "post-service.handler.postHandler.GetImageByImageID()"}).Error("Unable to encode image.")
		log.Println("unable to encode image.")
	}
	mediaZaFront = buffer.Bytes()
	imagesMarshaled, err := json.Marshal(mediaZaFront)
	if err != nil {
		log.WithFields(logrus.Fields{
			"location": "post-service.handler.postHandler.GetImageByImageID()"}).Error(err)
		fmt.Println(err)
	}
	log.WithFields(logrus.Fields{
		"location": "post-service.handler.postHandler.GetImageByImageID()"}).Info("Get image by imageID success.")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(imagesMarshaled)
}
