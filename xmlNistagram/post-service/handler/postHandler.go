package handler

import (
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
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

	path:=filepath.Join(absPath,"files",handle.Filename)
	f,err := os.OpenFile(path,os.O_WRONLY|os.O_CREATE,0666)

	if err != nil {
		http.Error(w, "Expected file", http.StatusBadRequest)
		return
	}
	io.Copy(f,file)
	jsonResponse(w, http.StatusCreated, "File uploaded successfully!.")


}

func saveFile(w http.ResponseWriter, file multipart.File, handle *multipart.FileHeader) {
	/*var path = "files/" + handle.Filename
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 777)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer f.Close()
	io.Copy(f, file)
	*/
	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
		fmt.Fprintf(w, "%v", err)
		return
	}

	err = ioutil.WriteFile("post-service/files"+handle.Filename, data, 0666)
	if err != nil {
		fmt.Println(err)
		fmt.Fprintf(w, "%v", err)
		return
	}
	jsonResponse(w, http.StatusCreated, "File uploaded successfully!.")
}
func jsonResponse(w http.ResponseWriter, code int, message string) {
	//	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	fmt.Fprint(w, message)
}
