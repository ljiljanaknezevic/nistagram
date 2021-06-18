package handler

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"story-service-mod/model"
	"story-service-mod/service"

	"github.com/gorilla/mux"
	logrus "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

type StoryHandler struct {
	Service *service.StoryService
}

var log = logrus.New()

func init() {
	absPath, err := os.Getwd()

	path := filepath.Join(absPath, "files", "story-service.log")
	filel, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.Out = filel
	} else {
		log.WithFields(
			logrus.Fields{
				"location": "story-service.handler.storyHandler.init()",
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
			"location": "story-service.handler.storyHandler.init()",
		},
	).Info("Story-service Log file created")
}
func (handler *StoryHandler) SaveStory(w http.ResponseWriter, r *http.Request) {
	log.WithFields(logrus.Fields{
		"location":   "story-service.handler.storyHandler.SaveStory()",
		"user_email": template.HTMLEscapeString(r.PostFormValue("email"))}).Info("User add story.")

	r.ParseMultipartForm(32 << 20)
	file, handle, err := r.FormFile("file")

	if err != nil {
		log.WithFields(
			logrus.Fields{
				"location": "story-service.handler.storyHandler.SaveStory()",
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
				"location": "story-service.handler.storyHandler.SaveStory()",
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
			"location": "story-service.handler.storyHandler.SaveStory()"}).Error("Failed in creating file.")
		var err model.Error
		err = model.SetError(err, "Failed in creating file.")
		json.NewEncoder(w).Encode(err)
		w.WriteHeader(http.StatusExpectationFailed)
		return
	}
	log.WithFields(logrus.Fields{
		"location": "story-service.handler.storyHandler.SaveStory()"}).Info("File localy created.")
	fileId := handler.Service.FindFileIdByPath(path)

	var story model.Story
	story.Description = r.PostFormValue("description")
	story.Location = r.PostFormValue("location")
	story.Tags = r.PostFormValue("tags")
	story.ImageID = fileId
	story.Email = r.PostFormValue("email")
	err = handler.Service.SaveStory(&story)
	if err != nil {
		log.WithFields(logrus.Fields{
			"location": "story-service.handler.storyHandler.SaveStory()"}).Error("Save post failed.")
		var err model.Error
		err = model.SetError(err, "Save post fail.")
		json.NewEncoder(w).Encode(err)
		w.WriteHeader(http.StatusExpectationFailed)
		return
	}
	log.WithFields(logrus.Fields{
		"location": "story-service.handler.storyHandler.SaveStory()"}).Info("Save post success.")
	log.WithFields(logrus.Fields{
		"location":   "story-service.handler.storyHandler.SaveStory()",
		"user_email": template.HTMLEscapeString(r.PostFormValue("email"))}).Info("User add post success.")

	jsonResponse(w, http.StatusCreated, "File uploaded successfully!.")

}

func jsonResponse(w http.ResponseWriter, code int, message string) {
	//	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	fmt.Fprint(w, message)
}

func (handler *StoryHandler) GetAllStoriesByEmail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	email := vars["email"]
	var result []model.Story
	result = handler.Service.GetAllStoriesByEmail(email)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}



func (handler *StoryHandler) GetImageByImageID(w http.ResponseWriter, r *http.Request) {
	/*	vars := mux.Vars(r)
			imageID := vars["imageID"]
			fmt.Println(imageID)
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
		w.Write(imagesMarshaled)*/
}
